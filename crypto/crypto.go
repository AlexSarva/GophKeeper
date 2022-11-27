package crypto

type Crypto interface {
	Verify(payload string, signature64 string) bool
	Sign(payload string) (string, error)
	Encrypt(payload string) (string, error)
	Decrypt(payload string) (string, error)
}

type Cryptorizer struct {
	Cryptorizer Crypto
}
