package cryptorizer

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
)

type Cryptorizer struct {
	aesgcm cipher.AEAD
	nonce  []byte
}

func BestCryptorizer(secret string) *Cryptorizer {
	key := sha256.Sum256([]byte(secret))
	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	nonce := key[len(key)-aesgcm.NonceSize():]
	return &Cryptorizer{
		aesgcm: aesgcm,
		nonce:  nonce,
	}
}

func (c *Cryptorizer) Encrypt(text string) string {
	src := []byte(text)
	dst := c.aesgcm.Seal(nil, c.nonce, src, nil)
	return hex.EncodeToString(dst)
}

func (c *Cryptorizer) Decrypt(text string) (string, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return "", err
	}
	res, err := c.aesgcm.Open(nil, c.nonce, data, nil) // расшифровываем
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}
	return string(res), nil
}
