package server

import (
	"errors"
	"fmt"

	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/wire"
)

func authenticate(r respondent, ctx *context) (e error) {
	var authRequest *wire.AuthenticationRequest
	if authRequest, e = r.getInitialMessage(); e != nil {
		return
	}

	if authRequest.MethodName != wire.AuthenticationMethodPassword && authRequest.MethodName != wire.AuthenticationMethodPublicKey {
		return errors.New(fmt.Sprint("Unknown auth request method - ", authRequest.MethodName))
	}

	var marshaledPublicKey []byte
	if marshaledPublicKey, e = common.GetMarshalPublicKey(ctx.flags.PrivateKeyPath); e != nil {
		return
	}

	authResponse := &wire.AuthenticationResponse{
		MethodName: authRequest.MethodName,
		PublicKey:  marshaledPublicKey,
	}

	e = r.respondMessage(authResponse)
	return
}
