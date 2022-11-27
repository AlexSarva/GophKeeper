package utils

import (
	"AlexSarva/GophKeeper/utils/luhn"
	"regexp"
	"strconv"
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
