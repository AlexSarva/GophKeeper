package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
)

// ErrNotValidSing - error that occurs when it is impossible to recover UserID
var ErrNotValidSing = errors.New("sign is not valid")

// Encrypt convert user uuid to hash
// secret key should be the same for Encrypt and Decrypt
func Encrypt(data string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data)[:])
	dst := h.Sum(nil)
	var fullData []byte
	fullData = append(fullData, []byte(data)[:]...)
	fullData = append(fullData, dst...)
	return hex.EncodeToString(fullData)
}

// Decrypt convert user hash to uuid
func Decrypt(hashString string, secret []byte) (string, error) {
	var (
		data []byte // декодированное сообщение с подписью
		err  error
		sign []byte // HMAC-подпись от идентификатора
	)

	data, err = hex.DecodeString(hashString)
	if err != nil {
		log.Println(err)
		return "", ErrNotValidSing
	}
	id := string(data[:16])
	log.Println(id)
	h := hmac.New(sha256.New, secret)
	h.Write(data[:16])
	sign = h.Sum(nil)

	if hmac.Equal(sign, data[16:]) {
		return id, nil
	} else {
		return "", ErrNotValidSing
	}
}
