package client

import (
	"errors"

	"github.com/murphybytes/ucp/wire"
)

func auth(ctx *context, passwdReader func(a ...interface{}) (i int, e error)) (e error) {

	var authResponse *wire.AutenticationResponse
	if authResponse, e = ctx.server.initializeSecureChannel(); e != nil {
		return
	}

	if authResponse.AllowedAuthenticationMethod == wire.AuthenticationMethodPassword {
		// get password and authenticate
	} else if authResponse.AllowedAuthenticationMethod != wire.AuthenticationMethodPublicKey {
		// wot!?
		return errors.New("Server sent back unexpected method")
	}

	return
}
