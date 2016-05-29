package client

import (
	"bytes"
	"crypto/rsa"
	"encoding/gob"
	"math/big"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/wire"
)

type requester interface {
	initializeSecureChannel() (e error)
	get(request proto.Message) (response proto.Message, e error)
	close()
}

type server struct {
	conn       net.Conn
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func newServer(conn net.Conn, ctx *context) (r requester, e error) {

	var privateKey *rsa.PrivateKey
	if privateKey, e = common.GetPrivateKey(ctx.flags.PrivateKeyPath); e != nil {
		return
	}

	r = &server{
		conn:       conn,
		privateKey: privateKey,
	}
	return
}

// All we are doing here is exchanging public keys
func (s *server) initializeSecureChannel() (e error) {

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if e = encoder.Encode(s.privateKey.PublicKey); e != nil {
		return
	}

	if _, e = s.conn.Write(buffer.Bytes()); e != nil {
		return
	}

	response := make([]byte, wire.ReadBufferSize)
	if _, e = s.conn.Read(response); e != nil {
		return
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(response))

	s.publicKey = &rsa.PublicKey{
		N: &big.Int{},
	}

	e = decoder.Decode(s.publicKey)

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
