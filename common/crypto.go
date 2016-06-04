package common

import (
	"bytes"
	"crypto"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/gob"
	"encoding/pem"
	"errors"
	"hash"
	"io/ioutil"
	"os"
)

// KeyBufferFetcher returns an array of bytes containing a crypto key
type KeyBufferFetcher func(*Flags) ([]byte, error)

// generates public/private keys and write each to file
func ucpKeyGenerate(privateKeyPath, publicKeyPath string) (e error) {
	var privateKey *rsa.PrivateKey

	if privateKey, e = rsa.GenerateKey(rand.Reader, KeySize); e != nil {
		return
	}

	var publicKey rsa.PublicKey
	publicKey = privateKey.PublicKey

	var pemkey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	var privateKeyFile *os.File
	if privateKeyFile, e = os.Create(privateKeyPath); e != nil {
		return
	}
	defer privateKeyFile.Close()

	if e = pem.Encode(privateKeyFile, pemkey); e != nil {
		return
	}

	var encodedKeyBuffer []byte
	if encodedKeyBuffer, e = CreateBase64EncodedPublicKey(publicKey); e != nil {
		return
	}

	if e = ioutil.WriteFile(publicKeyPath, encodedKeyBuffer, 0655); e != nil {
		return
	}

	return

}

// CreateBase64EncodedPublicKey returns a textual representation of the pubilc
// key suitable for authorized_keys files
func CreateBase64EncodedPublicKey(key rsa.PublicKey) (encodedKey []byte, e error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if e = encoder.Encode(key); e != nil {
		return
	}

	encodedKey = make([]byte, base64.StdEncoding.EncodedLen(buffer.Len()))
	base64.StdEncoding.Encode(encodedKey, buffer.Bytes())
	encodedKey = append(encodedKey, '\n')

	return
}

// GetPrivateKey returns a private key
func GetPrivateKey(privateKeyPath string) (key *rsa.PrivateKey, e error) {

	var buff []byte
	if buff, e = ioutil.ReadFile(privateKeyPath); e != nil {
		return
	}

	block, _ := pem.Decode(buff)

	if key, e = x509.ParsePKCS1PrivateKey(block.Bytes); e != nil {
		return
	}

	return

}

// EncryptOAEP encrypts a buffer
func EncryptOAEP(publicKey crypto.PublicKey, unencrypted []byte) (encrypted []byte, e error) {
	var md5Hash hash.Hash
	var label []byte
	md5Hash = md5.New()

	if key, ok := publicKey.(*rsa.PublicKey); ok {
		encrypted, e = rsa.EncryptOAEP(md5Hash, rand.Reader, key, unencrypted, label)
	} else {
		e = errors.New("Could not produce public key")
	}
	return
}

// DecryptOAEP decrypts a buffer
func DecryptOAEP(privateKey crypto.PrivateKey, encrypted []byte) (decrypted []byte, e error) {
	var md5Hash hash.Hash
	var label []byte
	md5Hash = md5.New()

	if key, ok := privateKey.(*rsa.PrivateKey); ok {
		decrypted, e = rsa.DecryptOAEP(md5Hash, rand.Reader, key, encrypted, label)
	} else {
		e = errors.New("Unable to produce private key")
	}

	return
}

// EncryptAES Encrypt a string with symmetric encryption
func EncryptAES(block cipher.Block, iv []byte, unencrypted []byte) (encrypted []byte) {

	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypted = make([]byte, len(unencrypted))
	encrypter.XORKeyStream(encrypted, unencrypted)
	return encrypted
}

// DecryptAES decrypts a a string with symmetric encryption
func DecryptAES(block cipher.Block, iv []byte, encrypted []byte) (unencrypted []byte) {

	decrypter := cipher.NewCFBDecrypter(block, iv)
	unencrypted = make([]byte, len(encrypted))
	decrypter.XORKeyStream(unencrypted, encrypted)
	return unencrypted

}
