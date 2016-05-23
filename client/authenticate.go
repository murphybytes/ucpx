package client

import (
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/wire"
)

func authenticate(ctx *context) (e error) {
	var publicKeyBuffer []byte
	if publicKeyBuffer, e = getPublicKey(ctx.flags); e != nil {
		ctx.logger.LogError("Unable to fetch public key - ", e.Error())
	}

	authRequest := &wire.AuthenticationRequest{
		UserName:   ctx.fileInfo.user,
		MethodName: wire.AuthenticationMethodPublicKey,
		PublicKey:  publicKeyBuffer,
	}

	var outBuffer []byte
	if outBuffer, e = proto.Marshal(authRequest); e != nil {
		return
	}

	if _, e = ctx.conn.Write(outBuffer); e != nil {
		return
	}

	//TODO: handle response

	return
}

func getPrivateKey(flags *common.Flags) (key []byte, e error) {
	key, e = ioutil.ReadFile(flags.PrivateKeyPath)
	return
}

func getPublicKey(flags *common.Flags) (key []byte, e error) {
	key, e = ioutil.ReadFile(flags.PublicKeyPath)
	return
}
