package client

import "github.com/murphybytes/ucp/wire"

func authenticate(ctx *context, passwdReader func(a ...interface{}) (i int, e error)) (e error) {

	var authResponse *wire.AuthenticationResponse
	if authResponse, e = ctx.server.initializeSecureChannel(); e != nil {
		return
	}

	if authResponse.MethodName == wire.AuthenticationMethodPublicKey {

	} else {
		// If we get here we need a password
		var password string
		if _, e = passwdReader(&password); e != nil {
			return
		}
	}

	// send password over secure channel
	//	if response, e = ctx.server.secureGet(request proto.Message)

	return
}
