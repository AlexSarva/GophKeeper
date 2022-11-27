package crypto

import "AlexSarva/GophKeeper/crypto/cryptorsa"

// Crypto interface that used for different types of crypt
type Crypto interface {
	Verify(payload string, signature64 string) bool
	Sign(payload string) (string, error)
	Encrypt(payload string) (string, error)
	Decrypt(payload string) (string, error)
}

// Cryptorizer used for implements different types of crypt
type Cryptorizer struct {
	Cryptorizer Crypto
}

// InitCryptorizer initializer of Cryptorizer struct
func InitCryptorizer(ketsPath string, size int) (*Cryptorizer, error) {
	cryptorizer := cryptorsa.InitRSACrypt(ketsPath, size)
	if initErr := cryptorizer.InitCrypto(); initErr != nil {
		return nil, initErr
	}
	return &Cryptorizer{
		Cryptorizer: cryptorizer,
	}, nil
}
