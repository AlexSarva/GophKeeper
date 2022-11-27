package symmetric

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"log"
)

type SymmetricCrypto struct {
	aesgcm cipher.AEAD
	nonce  []byte
}

func SymmCrypto(secret string) *SymmetricCrypto {
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
	return &SymmetricCrypto{
		aesgcm: aesgcm,
		nonce:  nonce,
	}
}

func (sc *SymmetricCrypto) Encrypt(payload []byte) []byte {
	dst := sc.aesgcm.Seal(nil, sc.nonce, payload, nil)
	return dst
}

func (sc *SymmetricCrypto) Decrypt(text []byte) ([]byte, error) {
	res, err := sc.aesgcm.Open(nil, sc.nonce, text, nil) // расшифровываем
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}
	return res, nil
}
