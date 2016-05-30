package client

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/gob"
	"math/big"

	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/udt.go/udt"
)

type requestorTestConn struct {
	udt.Conn
	read             []byte
	write            []byte
	clientPublicKey  *rsa.PublicKey
	serverPublicKey  *rsa.PublicKey
	serverPrivateKey *rsa.PrivateKey
}

func (r *requestorTestConn) Read(b []byte) (n int, e error) {

	r.serverPrivateKey, _ = rsa.GenerateKey(rand.Reader, common.KeySize)
	//r.serverPublicKey = &r.serverPrivateKey.PublicKey
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if e = encoder.Encode(&r.serverPrivateKey.PublicKey); e != nil {
		return
	}

	copy(b, buffer.Bytes())

	return
}

func (r *requestorTestConn) Write(b []byte) (n int, e error) {
	buffer := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buffer)
	r.clientPublicKey = &rsa.PublicKey{
		N: &big.Int{},
	}
	e = decoder.Decode(r.clientPublicKey)
	return
}

func newTestRequestor() (r requester) {

	privateKey, _ := rsa.GenerateKey(rand.Reader, common.KeySize)

	r = &server{
		privateKey: privateKey,
		conn:       &requestorTestConn{},
	}

	return
}

// func TestInitializeSecureChannel(t *testing.T) {
// 	var e error
// 	req := newTestRequestor()
// 	if e = req.initializeSecureChannel(); e != nil {
// 		t.Fatal("expected success ", e.Error())
// 	}
//
// 	testMsg := "this is a test message"
// 	var encrypted, decrypted []byte
// 	srv := req.(*server)
// 	conn := srv.conn.(*requestorTestConn)
// 	// client to server
// 	if encrypted, e = common.EncryptOAEP(srv.publicKey, []byte(testMsg)); e != nil {
// 		t.Fatal("Client encryption failed ", e.Error())
// 	}
// 	if decrypted, e = common.DecryptOAEP(conn.serverPrivateKey, encrypted); e != nil {
// 		t.Fatal("Server decryption failed ", e.Error())
// 	}
// 	if string(decrypted) != testMsg {
// 		t.Fatal("client to server decrypted string should match")
// 	}
//
// 	// server to client
// 	if encrypted, e = common.EncryptOAEP(conn.clientPublicKey, []byte(testMsg)); e != nil {
// 		t.Fatal("Server encryption failed")
// 	}
//
// 	if decrypted, e = common.DecryptOAEP(srv.privateKey, encrypted); e != nil {
// 		t.Fatal("Client decryption failed")
// 	}
//
// 	if string(decrypted) != testMsg {
// 		t.Fatal("server to client decrypted string should match")
// 	}
//
// }
