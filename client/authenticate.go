package client

import (
	"crypto"
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/wire"
	"golang.org/x/crypto/ssh"
)

func authenticate(ctx *context, passwdReader func(a ...interface{}) (i int, e error)) (e error) {
	// get our public key and send to server
	var publicKeyBuffer []byte
	if publicKeyBuffer, e = common.GetPublicKey(ctx.flags); e != nil {
		ctx.logger.LogError("Unable to fetch public key - ", e.Error())
	}

	authRequest := &wire.AuthenticationRequest{
		UserName:   ctx.fileInfo.user,
		MethodName: wire.AuthenticationMethodPublicKey,
		PublicKey:  publicKeyBuffer,
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
	var publicKey crypto.PublicKey
	if publicKey, e = ssh.ParsePublicKey(authResponse.PublicKey); e != nil {
		return
	}

	ctx.server.setPublicKey(publicKey)

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
