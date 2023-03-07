package utils

import (
	"AlexSarva/GophKeeper/utils/luhn"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"unicode"
)

// CheckCardNumber check available symbols for credit card number in GUI interface
func CheckCardNumber(text string) bool {
	var re = regexp.MustCompile(`[^\d\-\s]`)
	return !re.MatchString(text)
}

// CheckCardOwner check available symbols for credit card owner in GUI interface
func CheckCardOwner(text string) bool {
	var re = regexp.MustCompile(`[^a-zA-Z\.\s]`)
	return !re.MatchString(text)
}

// CheckCardExp check available symbols for credit card expired date in GUI interface
func CheckCardExp(text string) bool {
	var re = regexp.MustCompile(`[^\d\/]`)
	return !re.MatchString(text)
}

// CheckValidCardNumber check Luhn control sum for credit card number
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

type passwordCheck struct {
	number  bool
	upper   bool
	special bool
}

// PasswordChecker represents checker for user password
type PasswordChecker struct {
	length  int
	number  bool
	upper   bool
	special bool
	check   *passwordCheck
}

// InitPasswordChecker initializer of PasswordChecker struct
func InitPasswordChecker(length int, number, upper, special bool) *PasswordChecker {
	return &PasswordChecker{
		length:  length,
		number:  number,
		upper:   upper,
		special: special,
		check: &passwordCheck{
			number:  false,
			upper:   false,
			special: false,
		},
	}
}

// VerifyPassword checks all requirements for strong password
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
