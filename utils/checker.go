package utils

import (
	"AlexSarva/GophKeeper/utils/luhn"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"unicode"
)

func CheckCardNumber(text string) bool {
	var re = regexp.MustCompile(`[^\d\-\s]`)
	return !re.MatchString(text)
}

func CheckCardOwner(text string) bool {
	var re = regexp.MustCompile(`[^a-zA-Z\.\s]`)
	return !re.MatchString(text)
}

func CheckCardExp(text string) bool {
	var re = regexp.MustCompile(`[^\d\/]`)
	return !re.MatchString(text)
}

func CheckValidCardNumber(cardNum string) bool {
	intRegexp := regexp.MustCompile(`[^\d]`)
	digits := intRegexp.ReplaceAllString(cardNum, "")
	number, numberErr := strconv.Atoi(digits)
	if numberErr != nil {
		return false
	}
	if !luhn.Valid(number) {
		return false
	}
	return true
}

type PasswordCheck struct {
	number  bool
	upper   bool
	special bool
}

type PasswordChecker struct {
	length  int
	number  bool
	upper   bool
	special bool
	check   *PasswordCheck
}

func InitPasswordChecker(length int, number, upper, special bool) *PasswordChecker {
	return &PasswordChecker{
		length:  length,
		number:  number,
		upper:   upper,
		special: special,
		check: &PasswordCheck{
			number:  false,
			upper:   false,
			special: false,
		},
	}
}

func (pr *PasswordChecker) VerifyPassword(s string) error {
	letters := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			pr.check.number = true
			letters++
		case unicode.IsUpper(c):
			pr.check.upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			pr.check.special = true
			letters++
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
		}
	}
	if !(letters >= pr.length && pr.check.upper == pr.upper && pr.check.number == pr.number && pr.check.special == pr.special) {
		message := fmt.Sprintf("Not strong password! length: %d of min %d. numbers: %v. upper: %v. special: %v.",
			letters,
			pr.length,
			pr.check.number == pr.number,
			pr.check.upper == pr.upper,
			pr.check.special == pr.special)
		return errors.New(message)
	}
	return nil
}
