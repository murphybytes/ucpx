package client

import (
	"crypto"
	"crypto/rsa"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/wire"
	"golang.org/x/crypto/ssh"
)

type requester interface {
	initializeSecureChannel() (response *wire.AuthenticationResponse, e error)
	get(request proto.Message) (response proto.Message, e error)
	close()
}

type server struct {
	conn       net.Conn
	publicKey  crypto.PublicKey
	privateKey crypto.PrivateKey
}

func newServer(conn net.Conn, ctx *context) (r requester, e error) {

	var privateKey crypto.PrivateKey
	if privateKey, e = common.GetPrivateKey(ctx.flags.PrivateKeyPath); e != nil {
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

// All we are doing here is exchanging public keys
func (s *server) initializeSecureChannel() (response *wire.AuthenticationResponse, e error) {
	rsaPrivateKey := s.privateKey.(*rsa.PrivateKey)

	var sshPublicKey ssh.PublicKey
	if sshPublicKey, e = ssh.NewPublicKey(rsaPrivateKey.Public()); e != nil {
		return
	}

	// first message goes to server unencrypted
	request := &wire.AuthenticationRequest{
		MethodName: wire.AuthenticationMethodPublicKey,
		PublicKey:  sshPublicKey.Marshal(),
	}

	var requestBuffer []byte
	if requestBuffer, e = proto.Marshal(request); e != nil {
		return
	}

	if _, e = s.conn.Write(requestBuffer); e != nil {
		return
	}

	// response is encrypted
	responseBuffer := make([]byte, wire.ReadBufferSize)
	if _, e = s.conn.Read(responseBuffer); e != nil {
		return
	}

	var unencrypted []byte
	if unencrypted, e = common.DecryptOAEP(rsaPrivateKey, responseBuffer); e != nil {
		return
	}

	if e = proto.Unmarshal(unencrypted, response); e != nil {
		return
	}

	s.publicKey, e = ssh.ParsePublicKey(response.PublicKey)

	return
}

func (s *server) get(request proto.Message) (response proto.Message, e error) {
	var requestBuffer []byte

	if requestBuffer, e = proto.Marshal(request); e != nil {
		return
	}

	var encryptedRequestBuffer []byte
	if encryptedRequestBuffer, e = common.EncryptOAEP(s.publicKey, requestBuffer); e != nil {
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
	if responseBuffer, e = common.DecryptOAEP(s.privateKey, encryptedResponseBuffer); e != nil {
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
