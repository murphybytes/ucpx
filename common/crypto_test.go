package common

import (
	"crypto"
	"crypto/rsa"
	"fmt"
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

	var privateKey crypto.PrivateKey
	if privateKey, err = GetPrivateKey(privateKeyPath); err != nil {
		t.Fatal("Couldn't get private key -", err.Error())
	}

	original := "I am some unencrypted text."
	var encrypted []byte

	rsaPrivateKey := privateKey.(*rsa.PrivateKey)

	if encrypted, err = EncryptOAEP(rsaPrivateKey.Public(), []byte(original)); err != nil {
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
