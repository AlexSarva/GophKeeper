package cryptorsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// RSACrypt implements ID_RSA crypto methods
type RSACrypt struct {
	keysPath string
	idRsa    string
	idRsaPub string
	keySize  int
}

// InitRSACrypt initializer of RSACrypt struct
// path - folder path for store keys (id_rsa / id_rsa.pub)
// size - refer to the number of bits in a key used by a cryptographic algorithm
func InitRSACrypt(path string, size int) *RSACrypt {
	if size < 1024 {
		log.Fatalln("key size must by more or equal 1024")
	}
	return &RSACrypt{
		keysPath: path,
		idRsa:    filepath.Join(path, "id_rsa"),
		idRsaPub: filepath.Join(path, "id_rsa.pub"),
		keySize:  size,
	}
}

// InitCrypto checks for the presence of keys in the directory,
// if no keys are found - creates a pair personal and public keys
func (r *RSACrypt) InitCrypto() error {
	if _, err := os.Stat(r.keysPath); err != nil {
		log.Printf("keys will be add in this path: %s", r.keysPath)
		err = os.Mkdir(r.keysPath, 0700)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(r.idRsa); err != nil {
		// generate key pair
		// save private key
		// save public key
		keyPair, keyPairErr := r.generateKeyPair()
		if keyPairErr != nil {
			return keyPairErr
		}
		savePrivateErr := r.saveIDRsa(keyPair)
		if savePrivateErr != nil {
			return savePrivateErr
		}
		savePubErr := r.saveIDRsaPub(keyPair)
		if savePubErr != nil {
			return savePubErr
		}
	}

	return nil
}

func (r *RSACrypt) generateKeyPair() (*rsa.PrivateKey, error) {
	// generate key pair
	keyPair, err := rsa.GenerateKey(rand.Reader, r.keySize)
	if err != nil {
		return nil, err
	}

	// validate key
	err = keyPair.Validate()
	if err != nil {
		return nil, err
	}

	return keyPair, nil
}

func (r *RSACrypt) saveIDRsa(keyPair *rsa.PrivateKey) error {
	// private key stream
	privateKeyBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(keyPair),
	}

	// create file
	f, err := os.Create(r.idRsa)
	if err != nil {
		return err
	}

	err = pem.Encode(f, privateKeyBlock)
	if err != nil {
		return err
	}

	return nil
}

func (r *RSACrypt) saveIDRsaPub(keyPair *rsa.PrivateKey) error {
	// public key stream
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&keyPair.PublicKey)
	if err != nil {
		return err
	}

	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}

	// create file
	f, err := os.Create(r.idRsaPub)
	if err != nil {
		return err
	}

	err = pem.Encode(f, publicKeyBlock)
	if err != nil {
		return err
	}
	return nil
}

func (r *RSACrypt) getIDRsa() (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(r.idRsa)
	if err != nil {
		return nil, err
	}

	keyBlock, _ := pem.Decode(keyData)
	if keyBlock == nil {
		return nil, errors.New("fail get idrsa, invalid key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (r *RSACrypt) getIDRsaPub() (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(r.idRsaPub)
	if err != nil {
		return nil, err
	}

	keyBlock, _ := pem.Decode(keyData)
	if keyBlock == nil {
		return nil, errors.New("fail get id_rsa.pub, invalid key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(keyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	switch pk := publicKey.(type) {
	case *rsa.PublicKey:
		return pk, nil
	default:
		return nil, errors.New("fail get id_rsa.pub, invalid type")
	}
}

// Sign signs the file and returns the signature in base64
func (r *RSACrypt) Sign(payload string) (string, error) {
	// remove unwated characters and get sha256 hash of the payload
	replacer := strings.NewReplacer("\n", "", "\r", "", " ", "")
	msg := strings.TrimSpace(strings.ToLower(replacer.Replace(payload)))
	hashed := sha256.Sum256([]byte(msg))

	privateKey, privateKeyErr := r.getIDRsa()
	if privateKeyErr != nil {
		return "", privateKeyErr
	}

	// sign the hased payload
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	// return base64 encoded string
	return base64.StdEncoding.EncodeToString(signature), nil
}

// Verify check sign on payload
func (r *RSACrypt) Verify(payload string, signature64 string) bool {
	// decode base64 encoded signature
	signature, err := base64.StdEncoding.DecodeString(signature64)
	if err != nil {
		log.Println("ERROR: fail to base64 decode, ", err.Error())
		return false
	}

	// remove unwated characters and get sha256 hash of the payload
	replacer := strings.NewReplacer("\n", "", "\r", "", " ", "")
	msg := strings.TrimSpace(strings.ToLower(replacer.Replace(payload)))
	hashed := sha256.Sum256([]byte(msg))

	publicKey, publicKeyErr := r.getIDRsaPub()
	if publicKeyErr != nil {
		log.Println(publicKeyErr)
		return false
	}

	if verifyErr := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature); verifyErr != nil {
		log.Println(verifyErr)
		return false
	}

	return true
}

// Encrypt cipher payload via personal private key
func (r *RSACrypt) Encrypt(payload string) (string, error) {
	// params
	label := []byte("OAEP Encrypted")
	rnd := rand.Reader
	hash := sha256.New()

	publicKey, publicKeyErr := r.getIDRsaPub()
	if publicKeyErr != nil {
		return "", publicKeyErr
	}

	// encrypt with OAEP
	cipherText, err := rsa.EncryptOAEP(hash, rnd, publicKey, []byte(payload), label)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt deciphers payload via personal public key
func (r *RSACrypt) Decrypt(payload string) (string, error) {
	// decode base64 encoded signature
	label := []byte("OAEP Encrypted")
	msg, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", err
	}

	// params
	rnd := rand.Reader
	hash := sha256.New()

	privateKey, privateKeyErr := r.getIDRsa()
	if privateKeyErr != nil {
		return "", privateKeyErr
	}

	// decrypt with OAEP
	plainText, err := rsa.DecryptOAEP(hash, rnd, privateKey, msg, label)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
