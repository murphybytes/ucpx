package client

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"hash"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/wire"
	"golang.org/x/crypto/ssh"
)

type requester interface {
	get(request proto.Message) (response proto.Message, e error)
	secureGet(request proto.Message) (response proto.Message, e error)
	setPublicKey(publicKey crypto.PublicKey)
	close()
}

type server struct {
	conn       net.Conn
	publicKey  crypto.PublicKey
	privateKey crypto.PrivateKey
}

func newServer(conn net.Conn, ctx *context) (r requester, e error) {
	var buffer []byte
	if buffer, e = common.GetPrivateKey(ctx.flags); e != nil {
		return
	}

	var rawPrivateKey interface{}
	if rawPrivateKey, e = ssh.ParseRawPrivateKey(buffer); e != nil {
		return
	}

	var privateKey crypto.PrivateKey
	var ok bool
	if privateKey, ok = rawPrivateKey.(crypto.PrivateKey); !ok {
		e = errors.New("Could not create private key from pem file")
		return
	}

	r = &server{
		conn:       conn,
		privateKey: privateKey,
	}
	return
}

func (s *server) setPublicKey(publicKey crypto.PublicKey) {
	s.publicKey = publicKey
}

func (s *server) get(request proto.Message) (response proto.Message, e error) {
	var requestBuffer []byte

	if requestBuffer, e = proto.Marshal(request); e != nil {
		return
	}

	if _, e = s.conn.Write(requestBuffer); e != nil {
		return
	}

	responseBuffer := make([]byte, wire.ReadBufferSize)
	if _, e = s.conn.Read(responseBuffer); e != nil {
		return
	}

	if e = proto.Unmarshal(responseBuffer, response); e != nil {
		return
	}

	return
}

func (s *server) secureGet(request proto.Message) (response proto.Message, e error) {
	var requestBuffer []byte

	if requestBuffer, e = proto.Marshal(request); e != nil {
		return
	}

	var encryptedRequestBuffer []byte
	if encryptedRequestBuffer, e = encryptOAEP(s.publicKey, requestBuffer); e != nil {
		return
	}

	if _, e = s.conn.Write(encryptedRequestBuffer); e != nil {
		return
	}

	encryptedResponseBuffer := make([]byte, wire.ReadBufferSize)
	if _, e = s.conn.Read(encryptedResponseBuffer); e != nil {
		return
	}

	var responseBuffer []byte
	if responseBuffer, e = decryptOAEP(s.privateKey, encryptedResponseBuffer); e != nil {
		return
	}

	if e = proto.Unmarshal(responseBuffer, response); e != nil {
		return
	}

	return

}

func (s *server) close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func encryptOAEP(publicKey crypto.PublicKey, unencrypted []byte) (encrypted []byte, e error) {
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

func decryptOAEP(privateKey crypto.PrivateKey, encrypted []byte) (decrypted []byte, e error) {
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
