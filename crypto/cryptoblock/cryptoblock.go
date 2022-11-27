package cryptoblock

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"log"
)

// AEADCrypto implements Authenticated Encryption with Associated Data crypto methods
type AEADCrypto struct {
	aesgcm cipher.AEAD
	nonce  []byte
}

// InitAEADCrypto initializer of AEADCrypto struct
// secret - personal secret word
func InitAEADCrypto(secret string) *AEADCrypto {
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
	return &AEADCrypto{
		aesgcm: aesgcm,
		nonce:  nonce,
	}
}

// Encrypt cipher payload with AEAD and secret key
func (sc *AEADCrypto) Encrypt(payload []byte) []byte {
	dst := sc.aesgcm.Seal(nil, sc.nonce, payload, nil)
	return dst
}

// Decrypt deciphers payload with AEAD and secret key
func (sc *AEADCrypto) Decrypt(text []byte) ([]byte, error) {
	res, err := sc.aesgcm.Open(nil, sc.nonce, text, nil) // расшифровываем
	if err != nil {
		return nil, err
	}
	return res, nil
}
