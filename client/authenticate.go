package client

import (
	"crypto"
	"crypto/rsa"
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/wire"
	"golang.org/x/crypto/ssh"
)

func authenticate(ctx *context, passwdReader func(a ...interface{}) (i int, e error)) (e error) {

	var privateKey crypto.PrivateKey
	if privateKey, e = common.GetPrivateKey(ctx.flags.PrivateKeyPath); e != nil {
		return
	}

	rsaKey := privateKey.(*rsa.PrivateKey)

	var publicKey ssh.PublicKey
	if publicKey, e = ssh.NewPublicKey(rsaKey.PublicKey); e != nil {
		return
	}

	authRequest := &wire.AuthenticationRequest{
		UserName:   ctx.fileInfo.user,
		MethodName: wire.AuthenticationMethodPublicKey,
		PublicKey:  publicKey.Marshal(),
	}

	var response proto.Message
	if response, e = ctx.server.get(authRequest); e != nil {
		return
	}

	authResponse, ok := response.(*wire.AuthenticationResponse)
	if !ok {
		e = errors.New("Unexpected response from server")
	}

	// get public key server sent us to encrypt messages sent to server
	var serverKey crypto.PublicKey
	if serverKey, e = ssh.ParsePublicKey(authResponse.PublicKey); e != nil {
		return
	}

	ctx.server.setPublicKey(serverKey)

	// If server sends us PUBLIC_KEY the public key we sent was
	// found in the target users authorized_keys file so we're done
	if authResponse.MethodName == wire.AuthenticationMethodPublicKey {
		return
	}

	// If we get here we need a password
	var password string
	if _, e = passwdReader(&password); e != nil {
		return
	}
	// send password over secure channel
	//	if response, e = ctx.server.secureGet(request proto.Message)

	return
}
