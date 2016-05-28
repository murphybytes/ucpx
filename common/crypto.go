package common

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"hash"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

// KeyBufferFetcher returns an array of bytes containing a crypto key
type KeyBufferFetcher func(*Flags) ([]byte, error)

// generates public/private keys and write each to file
func ucpKeyGenerate(privateKeyPath, publicKeyPath string) (e error) {
	var privateKey *rsa.PrivateKey

	if privateKey, e = rsa.GenerateKey(rand.Reader, KeySize); e != nil {
		return
	}

	var publicKey *rsa.PublicKey
	publicKey = &privateKey.PublicKey

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

	var sshPublicKey ssh.PublicKey
	if sshPublicKey, e = ssh.NewPublicKey(publicKey); e != nil {
		return
	}

	if e = ioutil.WriteFile(publicKeyPath, ssh.MarshalAuthorizedKey(sshPublicKey), 0655); e != nil {
		return
	}

	return

}

// GetPrivateKey returns a private key
func GetPrivateKey(privateKeyPath string) (key crypto.PrivateKey, e error) {

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

// GetMarshalPublicKey gets public key in wire format
func GetMarshalPublicKey(privateKeyPath string) (keyBuff []byte, e error) {
	var cpk crypto.PrivateKey
	if cpk, e = GetPrivateKey(privateKeyPath); e != nil {
		return
	}

	var rsaPrivateKey *rsa.PrivateKey
	rsaPrivateKey = cpk.(*rsa.PrivateKey)

	var sshPublicKey ssh.PublicKey
	if sshPublicKey, e = ssh.NewPublicKey(rsaPrivateKey.Public()); e != nil {
		return
	}

	keyBuff = sshPublicKey.Marshal()
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
