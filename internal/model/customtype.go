package model

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type JSONDate time.Time

const dateFormat = "02.01.2006"

// MarshalJSON serijalizuje datum u JSON formatu
func (d JSONDate) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", time.Time(d).Format(dateFormat))
	return []byte(formatted), nil
}

// UnmarshalJSON deserijalizuje datum iz JSON formata
func (d *JSONDate) UnmarshalJSON(data []byte) error {
	strInput := strings.Trim(string(data), "\"")
	t, err := time.Parse(dateFormat, strInput)
	if err != nil {
		return errors.New("invalid date format")
	}
	*d = JSONDate(t)
	return nil
}

// String vraća datum kao string u definisanom formatu
func (d JSONDate) String() string {
	return time.Time(d).Format(dateFormat)
}
