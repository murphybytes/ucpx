package server

import (
	"bytes"
	"crypto/rsa"
	"encoding/gob"
	"math/big"

	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/wire"
)

type respondent interface {
	initializeSecureChannel() (e error)
	getMessage() ([]byte, error)
	sendMessage([]byte) (e error)
}

type client struct {
	clientKey *rsa.PublicKey
	serverKey *rsa.PrivateKey
	context   *context
}

func newClient(ctx *context) (r respondent, e error) {

	var serverKey *rsa.PrivateKey
	if serverKey, e = common.GetPrivateKey(ctx.flags.PrivateKeyPath); e != nil {
		return
	}

	r = &client{
		serverKey: serverKey,
		context:   ctx,
	}

	return

}

// first message from client is unencrypted and contains their public key
func (c *client) initializeSecureChannel() (e error) {
	c.context.logger.LogInfo("Beginning public key exchange with client")
	networkReadBuff := make([]byte, wire.ReadBufferSize)
	var readBytes int

	if readBytes, e = c.context.conn.Read(networkReadBuff); e != nil {
		return
	}

	encoderBuffer := bytes.NewBuffer(networkReadBuff[:readBytes])
	decoder := gob.NewDecoder(encoderBuffer)

	authRequest := wire.AuthenticationRequest{
		PublicKey: rsa.PublicKey{
			N: &big.Int{},
		},
	}

	if e = decoder.Decode(&authRequest); e != nil {
		return
	}

	c.context.logger.LogInfo("Successfully received authRequest for ", authRequest.UserName)

	// TODO: check authorization here

	c.clientKey = &authRequest.PublicKey

	// we now have clients public key, so send server public key to client

	authResponse := wire.AutenticationResponse{
		UserName:                    authRequest.UserName,
		PublicKey:                   c.serverKey.PublicKey,
		AllowedAuthenticationMethod: authRequest.RequestedAuthenticationMethod,
		Status:     wire.OK,
		StatusText: "OK",
	}

	encoderBuffer.Reset()
	encoder := gob.NewEncoder(encoderBuffer)

	if e = encoder.Encode(authResponse); e != nil {
		return
	}

	if _, e = c.context.conn.Write(encoderBuffer.Bytes()); e != nil {
		return
	}

	c.context.logger.LogInfo("Sent our public key to client")

	return
}

func (c *client) getMessage() (msg []byte, e error) {
	buffer := make([]byte, wire.ReadBufferSize)
	var read int
	if read, e = c.context.conn.Read(buffer); e != nil {
		return
	}

	if msg, e = common.DecryptOAEP(c.serverKey, buffer[:read]); e != nil {
		return
	}

	return msg, nil

}

func (c *client) sendMessage(msg []byte) (e error) {

	var encrypted []byte
	if encrypted, e = common.EncryptOAEP(c.clientKey, msg); e != nil {
		return
	}

	_, e = c.context.conn.Write(encrypted)

	return
}
