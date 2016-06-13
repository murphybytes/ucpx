package client

import (
	"bytes"
	"crypto/rsa"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"

	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/wire"
)

type requester interface {
	initializeSecureChannel() (*wire.AutenticationResponse, error)
	get([]byte) ([]byte, error)
	Write([]byte) (int, error)
	Read([]byte) (int, error)
	Close() error
}

type server struct {
	conn       net.Conn
	context    *context
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	buffer     []byte
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
		buffer:     make([]byte, wire.TxferBufferSize),
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

	s.publicKey = &authResponse.PublicKey

	return
}

func (s *server) get(request []byte) (response []byte, e error) {
	//	fmt.Printf("Public Key %q\n", s.publicKey)
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

func (s *server) Read(buff []byte) (n int, e error) {

	request := wire.ClientDataRequest{
		Status:     wire.More,
		StatusText: "More",
	}

	var encoderBuffer bytes.Buffer
	encoder := gob.NewEncoder(&encoderBuffer)

	if e = encoder.Encode(request); e != nil {
		return
	}

	encrypted := common.EncryptAES(s.context.aesKey, s.context.initializationVector, encoderBuffer.Bytes())

	if _, e = s.conn.Write(encrypted); e != nil {
		return
	}

	readBuffer := make([]byte, wire.ReadBufferSize)

	if n, e = s.conn.Read(readBuffer); e != nil {
		return
	}

	decrypted := common.DecryptAES(s.context.aesKey, s.context.initializationVector, readBuffer[:n])

	decodeBuffer := bytes.NewBuffer(decrypted)
	decoder := gob.NewDecoder(decodeBuffer)

	var clientDataResponse wire.ClientDataResponse

	if e = decoder.Decode(&clientDataResponse); e != nil {
		return
	}

	if clientDataResponse.Status == wire.EOF {
		e = io.EOF
		return
	}

	if clientDataResponse.Status != wire.OK {
		e = errors.New(clientDataResponse.StatusText)
		return
	}

	s.context.initializationVector = clientDataResponse.NextInitializationVector

	n = copy(buff, clientDataResponse.Data)

	return
}

func (s *server) Write(buff []byte) (n int, e error) {
	fmt.Printf("called write writing % data\n", len(buff))
	clientRead := &wire.ClientRead{
		Buffer:     buff,
		Status:     wire.More,
		StatusText: "More",
	}

	n = len(clientRead.Buffer)

	var encoderBuffer bytes.Buffer
	encoder := gob.NewEncoder(&encoderBuffer)
	if e = encoder.Encode(clientRead); e != nil {
		return
	}

	encrypted := common.EncryptAES(s.context.aesKey, s.context.initializationVector, encoderBuffer.Bytes())
	fmt.Printf("encrypted bytes %d\n", len(encrypted))
	var written int
	if written, e = s.conn.Write(encrypted); e != nil {
		return
	}

	fmt.Printf("wrote encoded %d bytes to server\n", written)

	readBuffer := make([]byte, wire.ReadBufferSize)

	var read int
	if read, e = s.conn.Read(readBuffer); e != nil {
		return
	}

	decrypted := common.DecryptAES(s.context.aesKey, s.context.initializationVector, readBuffer[:read])
	decodeBuffer := bytes.NewBuffer(decrypted)
	decoder := gob.NewDecoder(decodeBuffer)

	var clientReadResponse wire.ClientReadResponse
	if e = decoder.Decode(&clientReadResponse); e != nil {
		return
	}

	if clientReadResponse.Status != wire.OK {
		return 0, errors.New(clientReadResponse.StatusText)
	}

	s.context.initializationVector = clientReadResponse.NextInitializationVector

	return
}

func (s *server) Close() (e error) {
	if s.conn != nil {
		e = s.conn.Close()
	}
	return
}
