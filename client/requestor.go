package client

import (
	"crypto"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/wire"
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
