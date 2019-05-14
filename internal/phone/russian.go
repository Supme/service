package phone

import (
	"encoding/csv"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type russian struct {
	storage map[int][]line
	mu      sync.RWMutex
	update  time.Duration
}

func NewRussian(updateMin int) (*russian, error) {
	b := new(russian)
	bases := map[string]string{
		"ABC-3x.csv": "http://www.rossvyaz.ru/docs/articles/Kody_ABC-3kh.csv",
		"ABC-4x.csv": "http://www.rossvyaz.ru/docs/articles/Kody_ABC-4kh.csv",
		"ABC-8x.csv": "http://www.rossvyaz.ru/docs/articles/Kody_ABC-8kh.csv",
		"DEF-9x.csv": "http://www.rossvyaz.ru/docs/articles/Kody_DEF-9kh.csv",
	}
	b.update = time.Duration(updateMin) * time.Minute
	if b.update != 0 {
		err := b.updateRuBase(bases)
		if err != nil {
			return nil, err
		}
		go func(v *russian) {
			t := time.NewTicker(v.update)
			defer t.Stop()
			for range t.C {
				err := b.updateRuBase(bases)
				if err != nil {
					log.Printf("update russian base error: %s", err)
				}
			}
		}(b)
	} else {
		b.parseRuBase(bases)
	}
	return b, nil
}

func (b russian) Find(num []rune) (string, string, error) {
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

func (b *russian) parseRuBase(bases map[string]string) (map[int][]line, error) {
	storage := map[int][]line{}
	for k := range bases {
		f, err := os.Open(k)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		c := csv.NewReader(f)
		c.Comma = ';'
		c.LazyQuotes = true
		c.FieldsPerRecord = -1
		// skip first line
		_, err = c.Read()
		if err != nil {
			return nil, err
		}
		i := 1
		for {
			rec, err := c.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			i++
			if len(rec) < 5 {
				return nil, fmt.Errorf("file %s has wrong line %d", k, i)
				continue
			}
			f1, err := strconv.Atoi(rec[0])
			if err != nil {
				return nil, err
			}
			f2, err := strconv.Atoi(rec[1])
			if err != nil {
				return nil, err
			}
			f3, err := strconv.Atoi(rec[2])
			if err != nil {
				return nil, err
			}
			f4, err := strconv.Atoi(rec[3])
			if err != nil {
				return nil, err
			}
			if (f3 - f2 + 1) != f4 {
				return nil, fmt.Errorf("file %s has wrong count number in line %d (%d-%d+1) != %d", k, i, f3, f2, f4)
				continue
			}
			var f5 string
			if len(rec) == 6 {
				f5 = rec[4] + " " + rec[5]
			} else {
				f5 = rec[4]
			}
			storage[f1] = append(storage[f1], line{from: f2, to: f3, name: f5})
		}
	}

	// add Kazahstan
	for i := range kazahstan {
		for n := range kazahstan[i] {
			storage[i] = append(storage[i], kazahstan[i][n])
		}

	}
	return storage, nil
}

func (b *russian) updateRuBase(bases map[string]string) error {
	log.Print("start update russian base")
	for k, v := range bases {
		if fileInfo, err := os.Stat(k); os.IsNotExist(err) || fileInfo.ModTime().Add(b.update+(time.Second*59)).Unix() < time.Now().Unix() {
			//log.Printf("start download russian csv base %s\n", k)
			resp, err := http.Get(v)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if _, err := os.Stat(k); !os.IsNotExist(err) {
				os.Rename(k, "~"+k)
			}
			file, err := os.Create(k)
			if err != nil {
				if _, err := os.Stat("~" + k); !os.IsNotExist(err) {
					os.Rename("~"+k, k)
				}
				return err
			}
			defer file.Close()
			body := transform.NewReader(resp.Body, charmap.Windows1251.NewDecoder())
			_, err = io.Copy(file, body)
			if err != nil {
				if _, err := os.Stat("~" + k); !os.IsNotExist(err) {
					os.Rename("~"+k, k)
				}
				return err
			}
		}
	}

	storage, err := b.parseRuBase(bases)
	if err != nil {
		for k := range bases {
			if _, err := os.Stat("~" + k); !os.IsNotExist(err) {
				os.Remove(k)
				os.Rename("~"+k, k)
			}
		}
		return err
	}

	now := time.Now().Format("02-01-2006T15:04:05")
	for k := range bases {
		if _, err := os.Stat("~" + k); !os.IsNotExist(err) {
			os.Rename("~"+k, now+"_"+k)
		}
	}

	b.mu.Lock()
	b.storage = storage
	b.mu.Unlock()
	log.Print("russian base updated")
	return nil
}
