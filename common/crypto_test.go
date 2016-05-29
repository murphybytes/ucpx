package common

import (
	"bytes"
	"crypto/rsa"
	"encoding/gob"
	"fmt"
	"math/big"
	"os"
	"testing"
)

func CreateTestDirectory() (dir string, e error) {
	dir = fmt.Sprint(os.Getenv("GOPATH"), "/src/github.com/murphybytes/ucp/testdata")
	e = os.MkdirAll(dir, 0777)
	return
}

func DeleteTestDirectory(dir string) {
	os.RemoveAll(dir)
}

func TestEncryption(t *testing.T) {
	testdir, err := CreateTestDirectory()
	if err != nil {
		t.Fatal("Test data directory creation failed -", err.Error())
	}
	defer DeleteTestDirectory(testdir)

	publicKeyPath := fmt.Sprint(testdir, "/id_rsa.pub")
	privateKeyPath := fmt.Sprint(testdir, "/private.pem")

	if err = ucpKeyGenerate(privateKeyPath, publicKeyPath); err != nil {
		t.Fatal("Key generation failed -", err.Error())
	}

	var privateKey *rsa.PrivateKey
	if privateKey, err = GetPrivateKey(privateKeyPath); err != nil {
		t.Fatal("Couldn't get private key -", err.Error())
	}

	original := "I am some unencrypted text."
	var encrypted []byte

	if encrypted, err = EncryptOAEP(&privateKey.PublicKey, []byte(original)); err != nil {
		t.Fatal("EncryptOAEP failed -", err.Error())
	}

	if string(encrypted) == original {
		t.Fatal("Encrypted string should not match original")
	}

	var unencrypted []byte
	if unencrypted, err = DecryptOAEP(privateKey, encrypted); err != nil {
		t.Fatal("DecryptOAEP failed -", err.Error())
	}

	if string(unencrypted) != original {
		t.Fatal("Unencrypted string should match original")
	}

}

func TestEncryptionWithMarshalling(t *testing.T) {

	testdir, err := CreateTestDirectory()
	if err != nil {
		t.Fatal("Test data directory creation failed -", err.Error())
	}
	defer DeleteTestDirectory(testdir)

	publicKeyPath := fmt.Sprint(testdir, "/id_rsa.pub")
	privateKeyPath := fmt.Sprint(testdir, "/private.pem")

	if err = ucpKeyGenerate(privateKeyPath, publicKeyPath); err != nil {
		t.Fatal("Key generation failed -", err.Error())
	}

	var privateKey *rsa.PrivateKey
	if privateKey, err = GetPrivateKey(privateKeyPath); err != nil {
		t.Fatal("Couldn't get private key -", err.Error())
	}

	original := "I am some unencrypted text."
	var encrypted []byte

	var buffer bytes.Buffer
	decoder := gob.NewDecoder(&buffer)
	encoder := gob.NewEncoder(&buffer)

	if err = encoder.Encode(privateKey.PublicKey); err != nil {
		t.Fatal("Public key encode failed ", err.Error())
	}

	pk := &rsa.PublicKey{
		N: &big.Int{},
	}

	if err = decoder.Decode(pk); err != nil {
		t.Fatal("Public key decode failed - ", err.Error())
	}

	if encrypted, err = EncryptOAEP(pk, []byte(original)); err != nil {
		t.Fatal("EncryptOAEP failed -", err.Error())
	}

	if string(encrypted) == original {
		t.Fatal("Encrypted string should not match original")
	}

	var unencrypted []byte
	if unencrypted, err = DecryptOAEP(privateKey, encrypted); err != nil {
		t.Fatal("DecryptOAEP failed -", err.Error())
	}

	if string(unencrypted) != original {
		t.Fatal("Unencrypted string should match original")
	}

}
