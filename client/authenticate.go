package client

func authenticate(ctx *context, passwdReader func(a ...interface{}) (i int, e error)) (e error) {

	if e = ctx.server.initializeSecureChannel(); e != nil {
		return
	}

	// send password over secure channel
	//	if response, e = ctx.server.secureGet(request proto.Message)

	return
}
