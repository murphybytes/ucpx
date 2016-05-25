package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

// generates public/private keys and write each to file
func ucpKeyGenerate(privateKeyPath, publicKeyPath string) (e error) {
	var privateKey *rsa.PrivateKey

	if privateKey, e = rsa.GenerateKey(rand.Reader, KeySize); e != nil {
		return
	}

	var privateFile *os.File
	if privateFile, e = os.Create(privateKeyPath); e != nil {
		return
	}
	defer privateFile.Close()

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if e = pem.Encode(privateFile, privateKeyPEM); e != nil {
		return
	}

	var pub ssh.PublicKey
	if pub, e = ssh.NewPublicKey(&privateKey.PublicKey); e != nil {
		return
	}

	return ioutil.WriteFile(publicKeyPath, ssh.MarshalAuthorizedKey(pub), 0655)

}

// GetPrivateKey returns serialized private key file, file
// location may be specified on command line
func GetPrivateKey(flags *Flags) (key []byte, e error) {
	key, e = ioutil.ReadFile(flags.PrivateKeyPath)
	return
}

// GetPublicKey returns contents of public key file.  Non-default file
// location may be specified on command line
func GetPublicKey(flags *Flags) (key []byte, e error) {
	key, e = ioutil.ReadFile(flags.PublicKeyPath)
	return
}
