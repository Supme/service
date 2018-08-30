package phone

import (
	"fmt"
	"github.com/supme/service/proto"
	"golang.org/x/net/context"
	"io"
	"log"
	"unicode"
)

type Validator struct {
	ruBase base
}

type base interface {
	Find([]rune) (string, string, error)
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
	var err error
	validator.ruBase, err = NewRussian(43200)
	if err != nil {
		return nil, err
	}
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
	if len(num) == 10 && (num[0] != '+' || num[0] != '0') {
		return v.ruBase.Find(num)
	}
	// default prefix in Russia
	if num[0] == '8' {
		return v.ruBase.Find(num[1:])
	}

	if num[0] == '+' || num[0] == '0' {
		var phoneNum []rune
		// remove prefix
		if num[0] == '0' && num[1] == '0' {
			phoneNum = num[2:]
		} else {
			phoneNum = num[1:]
		}
		c := findCountry(phoneNum)
		switch c {
		// is Russia prefix?
		case "Россия":
			if len(phoneNum) > 11 {
				return "", "", errWrongLenghtNumber
			}
			return v.ruBase.Find(phoneNum[1:])
		case "":
			return "", "", errDontKnowCountryCode
		default:
			return "+" + string(phoneNum), c, nil
		}
	}
	return "", "", errDontKnowPhone
}
