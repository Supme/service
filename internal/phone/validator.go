package phone

import (
	"encoding/csv"
	"fmt"
	"github.com/supme/service/proto"
	"golang.org/x/net/context"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
	"unicode"
)

type Validator struct {
	ruBase base
}

type base struct {
	storage map[int][]line
	mu      sync.RWMutex
}

type line struct {
	from, to int
	name     string
}

var (
	errDontKnowCountryCode                         = fmt.Errorf("don't know country code")
	errDontKnowPhone                               = fmt.Errorf("don't know phone")
	errWrongLenghtNumber                           = fmt.Errorf("wrong lenght number")
	errCodeNotFoundForRussianDatabase              = fmt.Errorf("code not found for russian database")
	errNumberNotFoundInCodeRangeForRussianDatabase = fmt.Errorf("number not found in code range for russian database")
)

func NewValidator() (*Validator, error) {
	validator := new(Validator)
	err := validator.updateRuBase()
	if err != nil {
		return nil, err
	}
	go func(v *Validator) {
		t := time.NewTicker(12 * time.Hour)
		defer t.Stop()
		for range t.C {
			validator.updateRuBase()
		}
	}(validator)
	return validator, nil
}

func (v *Validator) StreamValidate(in proto.Phone_StreamValidateServer) error {
	for {
		r, err := in.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
		}
		s, err := v.Validate(in.Context(), r)
		err = in.Send(s)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (v *Validator) Validate(ctx context.Context, in *proto.PhoneValidateRequest) (*proto.PhoneValidateReply, error) {
	valid := false
	canonical, provider, err := v.check(in.Number)
	if err == nil {
		valid = true
	}

	resp := proto.PhoneValidateReply{
		Id:        in.Id,
		Valid:     valid,
		Canonical: canonical,
		Provider:  provider,
		Error:     phoneErrorToProtoError(err),
	}
	return &resp, nil
}

func phoneErrorToProtoError(err error) proto.PhoneValidateError {
	switch err {
	case nil:
		return proto.PhoneValidateError_NO_ERROR
	case errDontKnowCountryCode:
		return proto.PhoneValidateError_DONT_KNOW_COUNTRY_CODE
	case errDontKnowPhone:
		return proto.PhoneValidateError_DONT_KNOW_PHONE
	case errWrongLenghtNumber:
		return proto.PhoneValidateError_WRONG_LENGHT_NUMBER
	case errCodeNotFoundForRussianDatabase:
		return proto.PhoneValidateError_CODE_NOT_FOUND_FOR_RUSSIAN_DATABASE
	case errNumberNotFoundInCodeRangeForRussianDatabase:
		return proto.PhoneValidateError_NUMBER_NOT_FOUND_IN_CODE_RANGE_FOR_RUSSIAN_DATABASE
	default:
		return proto.PhoneValidateError_OTHER_ERROR
	}
}

// Check returns canonical phone format, provider or country for phone
func (v *Validator) Check(number string) (string, string, error) {
	return v.check(number)
}

func (v *Validator) check(number string) (string, string, error) {
	num := make([]rune, 0, 15)
	// clean number
	for i := range number {
		if unicode.IsDigit(rune(number[i])) || rune(number[i]) == '+' {
			num = append(num, rune(number[i]))
		}
	}
	if len(num) < 10 {
		return "", "", errWrongLenghtNumber
	}
	// Russia default country
	if len(num) == 10 {
		return v.ruBase.find(num)
	}
	// default prefix in Russia
	if num[0] == '8' {
		return v.ruBase.find(num[1:])
	}

	if num[0] == '+' || num[0] == '0' {
		var phoneNum []rune
		if num[0] == '0' && num[1] == '0' {
			phoneNum = num[2:]
		} else {
			phoneNum = num[1:]
		}
		c := findCountry(phoneNum)
		switch c {
		// is Russia prefix?
		case "Россия":
			return v.ruBase.find(phoneNum[1:])
		case "":
			return "", "", errDontKnowCountryCode
		default:
			return "+" + string(phoneNum), c, nil
		}
	}
	return "", "", errDontKnowPhone
}

func (b base) find(num []rune) (string, string, error) {
	if len(num) != 10 {
		return "", "", errWrongLenghtNumber
	}
	code, err := strconv.Atoi(string(num[0:3]))
	if err != nil {
		return "", "", err
	}
	number, err := strconv.Atoi(string(num[3:10]))
	if err != nil {
		return "", "", err
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	rec, ok := b.storage[code]
	if !ok {
		return "", "", errCodeNotFoundForRussianDatabase
	}
	var provider string
	found := false
	for i := range rec {
		if number >= rec[i].from && number <= rec[i].to {
			provider = rec[i].name
			found = true
			break
		}
	}
	if !found {
		return "", "", errNumberNotFoundInCodeRangeForRussianDatabase
	}
	return "+7" + string(num), provider, nil
}

func (v *Validator) updateRuBase() error {
	bases := map[string]string{
		"ABC-3x.csv": "http://www.rossvyaz.ru/docs/articles/ABC-3x.csv",
		"ABC-4x.csv": "http://www.rossvyaz.ru/docs/articles/ABC-4x.csv",
		"ABC-8x.csv": "http://www.rossvyaz.ru/docs/articles/ABC-8x.csv",
		"DEF-9x.csv": "http://www.rossvyaz.ru/docs/articles/DEF-9x.csv",
	}
	err := v.ruBase.updateRuBaseCSV(bases)
	if err != nil {
		if !v.ruBase.existRuBaseCSV(bases) {
			return err
		}
		log.Print("Update russian csv files this error, continue use old file")
	} else {
		err = v.ruBase.parseRuBase(bases)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *base) parseRuBase(bases map[string]string) error {
	log.Printf("Start parse russian csv base")
	b.mu.Lock()
	defer b.mu.Unlock()
	b.storage = map[int][]line{}
	for k := range bases {
		f, err := os.Open(k)
		if err != nil {
			return err
		}
		defer f.Close()
		c := csv.NewReader(f)
		c.Comma = ';'
		c.LazyQuotes = true
		c.FieldsPerRecord = -1
		// skip first line
		_, err = c.Read()
		if err != nil {
			return err
		}
		i := 1
		for {
			rec, err := c.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			i++
			if len(rec) < 5 {
				log.Printf("File %s has wrong line %d", k, i)
				continue
			}
			f1, err := strconv.Atoi(rec[0])
			if err != nil {
				return err
			}
			f2, err := strconv.Atoi(rec[1])
			if err != nil {
				return err
			}
			f3, err := strconv.Atoi(rec[2])
			if err != nil {
				return err
			}
			f4, err := strconv.Atoi(rec[3])
			if err != nil {
				return err
			}
			if (f3 - f2 + 1) != f4 {
				log.Printf("File %s has wrong count number in line %d (%d-%d+1) != %d", k, i, f3, f2, f4)
				continue
			}
			var f5 string
			if len(rec) == 6 {
				f5 = rec[4] + " " + rec[5]
			} else {
				f5 = rec[4]
			}
			b.storage[f1] = append(b.storage[f1], line{from: f2, to: f3, name: f5})
		}
	}

	// add Kazahstan
	for i := range kazahstan {
		for n := range kazahstan[i] {
			b.storage[i] = append(b.storage[i], kazahstan[i][n])
		}

	}
	return nil
}

func (b base) existRuBaseCSV(bases map[string]string) bool {
	exist := true
	for k := range bases {
		if _, err := os.Stat(k); os.IsNotExist(err) {
			exist = false
			break
		}
	}
	return exist
}

func (b base) updateRuBaseCSV(bases map[string]string) error {
	for k, v := range bases {
		if fileInfo, err := os.Stat(k); os.IsNotExist(err) || fileInfo.ModTime().Add(time.Hour*24).Unix() < time.Now().Unix() {
			log.Printf("Start download russian csv base %s\n", k)
			resp, err := http.Get(v)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			// ToDo use temp file for download and rename old
			file, err := os.Create(k)
			if err != nil {
				return err
			}
			defer file.Close()
			body := transform.NewReader(resp.Body, charmap.Windows1251.NewDecoder())
			_, err = io.Copy(file, body)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
