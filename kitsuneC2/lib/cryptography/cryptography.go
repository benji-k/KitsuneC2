//Package containing general cryptographic functions used both by the server and the implant. The library has nice wrappers for AES encryption/decryption
//and hashing functions.

package cryptography

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

// Generates a new RSA key pair. 2048 bits is considered the standard.
func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privkey, &privkey.PublicKey, nil
}

// RSA private key to base64 encoded string
func RsaPrivateKeyToString(priv *rsa.PrivateKey) string {
	privKeyBytes := x509.MarshalPKCS1PrivateKey(priv)
	return base64.StdEncoding.EncodeToString(privKeyBytes)
}

// Given a base64 encoded RSA private key, returns an *rsa.PrivateKey object
func StringToRsaPrivateKey(priv string) (*rsa.PrivateKey, error) {
	privKeyBytes, err := base64.StdEncoding.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKCS1PrivateKey(privKeyBytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// RSA public key to base64 encoded string
func RSAPublicKeyToString(pub *rsa.PublicKey) string {
	pubKeyBytes := x509.MarshalPKCS1PublicKey(pub)
	return base64.StdEncoding.EncodeToString(pubKeyBytes)
}

// Given a base64 encoded RSA public key, returns an *rsa.PublicKey object
func StringToRsaPublicKey(pub string) (*rsa.PublicKey, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(pub)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKCS1PublicKey(publicKeyBytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// Encrypts data with public RSA key
func EncryptWithRsaPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// Decrypts data with private RSA key
func DecryptWithRsaPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func GenerateMd5FromStrings(strings ...string) (md5Hash string) {
	var concatenated string = ""
	for _, s := range strings {
		concatenated += s
	}
	hash := md5.Sum([]byte(concatenated))
	return hex.EncodeToString(hash[:])
}

// plaintext needs to be padded to a multiple of aes.Blocksize(usually 16). We pad the plaintext with repeating bytes. One such byte represents
// the amount of added bytes for padding. e.g. If our blockksize - plaintext % blocksize = 2. This means we need to add 2 extra bytes. The value of these
// bytes will be 0x02, representing the amount of bytes we have added. Do note that if our plaintext is exactly a multiple of aes.blockSize, we append
// one whole block of padding to the plaintext.
func pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// By reading the last byte of our decrypted cipertext, we know how many bytes to remove from the plaintext to get back our original input.
func unPad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])
	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}
	return src[:(length - unpadding)], nil
}

// EncryptAes takes in a plaintext string and a key and returns the encrypted data. They key needs to be either 16, 24 or 32 bits.
func EncryptAes(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	data = pad(data)
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], data)
	return ciphertext, nil
}

// DecryptAes takes in a cipher text and decryptes the data. They key needs to be either 16, 24 or 32 bits.
func DecryptAes(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext, err = unPad(ciphertext)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}
