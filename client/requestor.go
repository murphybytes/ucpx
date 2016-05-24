package client

import (
	"crypto"
	"net"

	"github.com/golang/protobuf/proto"
)

type requester interface {
	get(request proto.Message) (response proto.Message, e error)
	secureGet(request proto.Message) (response proto.Message, e error)
	close()
}

type server struct {
	conn       net.Conn
	publicKey  crypto.PublicKey
	privateKey crypto.PrivateKey
}

func newServer(conn net.Conn) requester {
	return &server{
		conn: conn,
	}
}

func (s *server) get(request proto.Message) (response proto.Message, e error) {
	var requestBuffer []byte

	if requestBuffer, e = proto.Marshal(request); e != nil {
		return
	}

	if _, e = s.conn.Write(requestBuffer); e != nil {
		return
	}

	var responseBuffer []byte
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

	encryptedRequestBuffer := encryptOAEP(s.publicKey, requestBuffer)

	if _, e = s.conn.Write(encryptedRequestBuffer); e != nil {
		return
	}

	var encryptedResponseBuffer []byte
	if _, e = s.conn.Read(encryptedResponseBuffer); e != nil {
		return
	}

	responseBuffer := decryptOAEP(s.privateKey, encryptedResponseBuffer)

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

func encryptOAEP(publicKey crypto.PublicKey, unencrypted []byte) (encrypted []byte) {
	return
}

func decryptOAEP(privateKey crypto.PrivateKey, encrypted []byte) (descrypted []byte) {
	return
}
