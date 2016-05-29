package server

import (
	"bytes"
	"crypto/rsa"
	"encoding/gob"
	"math/big"

	"github.com/golang/protobuf/proto"
	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/wire"
)

type respondent interface {
	initializeSecureChannel() (e error)
	getMessage() (msg proto.Message, e error)
	respondMessage(msg proto.Message) (e error)
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

	if _, e = c.context.conn.Read(networkReadBuff); e != nil {
		return
	}

	encoderBuffer := bytes.NewBuffer(networkReadBuff)
	decoder := gob.NewDecoder(encoderBuffer)

	c.clientKey = &rsa.PublicKey{
		N: &big.Int{},
	}

	if e = decoder.Decode(c.clientKey); e != nil {
		return
	}

	c.context.logger.LogInfo("Successfully received public key from client")

	// we now have clients public key, so send server public key to client
	encoderBuffer.Reset()
	encoder := gob.NewEncoder(encoderBuffer)

	if e = encoder.Encode(c.serverKey.PublicKey); e != nil {
		return
	}

	if _, e = c.context.conn.Write(encoderBuffer.Bytes()); e != nil {
		return
	}

	c.context.logger.LogInfo("Sent our public key to client")

	return
}

func (c *client) getMessage() (msg proto.Message, e error) {
	buffer := make([]byte, wire.ReadBufferSize)
	if _, e = c.context.conn.Read(buffer); e != nil {
		return
	}

	var decrypted []byte

	if decrypted, e = common.DecryptOAEP(c.serverKey, buffer); e != nil {
		return
	}

	e = proto.Unmarshal(decrypted, msg)
	return

}

func (c *client) respondMessage(msg proto.Message) (e error) {
	var unencrypted []byte
	if unencrypted, e = proto.Marshal(msg); e != nil {
		return
	}

	var encrypted []byte
	if encrypted, e = common.EncryptOAEP(c.clientKey, unencrypted); e != nil {
		return
	}

	_, e = c.context.conn.Write(encrypted)

	return
}
