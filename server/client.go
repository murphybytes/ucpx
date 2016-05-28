package server

import (
	"crypto"

	"github.com/golang/protobuf/proto"
	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/wire"
	"golang.org/x/crypto/ssh"
)

type respondent interface {
	getInitialMessage() (clientAuthRequest *wire.AuthenticationRequest, e error)
	getMessage() (msg proto.Message, e error)
	respondMessage(msg proto.Message) (e error)
}

type client struct {
	clientKey crypto.PublicKey
	serverKey crypto.PrivateKey
	context   *context
}

func newClient(ctx *context) (r respondent, e error) {

	var serverKey crypto.PrivateKey
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
func (c *client) getInitialMessage() (clientAuthRequest *wire.AuthenticationRequest, e error) {
	buffer := make([]byte, wire.ReadBufferSize)
	if _, e = c.context.conn.Read(buffer); e != nil {
		return
	}

	if e = proto.Unmarshal(buffer, clientAuthRequest); e != nil {
		return
	}

	// set client public key that we will use to encrypt messages that go back to
	// client
	if c.clientKey, e = ssh.ParsePublicKey(clientAuthRequest.PublicKey); e != nil {
		return
	}

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
