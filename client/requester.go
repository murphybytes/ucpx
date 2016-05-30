package client

import (
	"bytes"
	"crypto/rsa"
	"encoding/gob"
	"errors"
	"math/big"
	"net"

	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/wire"
)

type requester interface {
	initializeSecureChannel() (*wire.AutenticationResponse, error)
	get([]byte) ([]byte, error)
	close()
}

type server struct {
	conn       net.Conn
	context    *context
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func newServer(conn net.Conn, ctx *context) (r requester, e error) {

	var privateKey *rsa.PrivateKey
	if privateKey, e = common.GetPrivateKey(ctx.flags.PrivateKeyPath); e != nil {
		return
	}

	r = &server{
		context:    ctx,
		conn:       conn,
		privateKey: privateKey,
	}
	return
}

// All we are doing here is exchanging public keys
func (s *server) initializeSecureChannel() (authResponse *wire.AutenticationResponse, e error) {

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	authRequest := &wire.AuthenticationRequest{
		UserName:                      s.context.fileInfo.user,
		RequestedAuthenticationMethod: wire.AuthenticationMethodPublicKey,
		PublicKey:                     s.privateKey.PublicKey,
	}

	if e = encoder.Encode(authRequest); e != nil {
		return
	}

	if _, e = s.conn.Write(buffer.Bytes()); e != nil {
		return
	}

	response := make([]byte, wire.ReadBufferSize)
	var read int
	if read, e = s.conn.Read(response); e != nil {
		return
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(response[:read]))

	authResponse = &wire.AutenticationResponse{
		PublicKey: rsa.PublicKey{
			N: &big.Int{},
		},
	}

	if e = decoder.Decode(authResponse); e != nil {
		return
	}

	if authResponse.Status != wire.OK {
		e = errors.New(authResponse.StatusText)
		return
	}

	return
}

func (s *server) get(request []byte) (response []byte, e error) {
	var encryptedRequestBuffer []byte
	if encryptedRequestBuffer, e = common.EncryptOAEP(s.publicKey, request); e != nil {
		return
	}

	if _, e = s.conn.Write(encryptedRequestBuffer); e != nil {
		return
	}

	encryptedResponseBuffer := make([]byte, wire.ReadBufferSize)
	var read int
	if read, e = s.conn.Read(encryptedResponseBuffer); e != nil {
		return
	}

	var responseBuffer []byte
	if responseBuffer, e = common.DecryptOAEP(s.privateKey, encryptedResponseBuffer[:read]); e != nil {
		return
	}

	return responseBuffer, nil

}

func (s *server) close() {
	if s.conn != nil {
		s.conn.Close()
	}
}
