package server

import (
	"crypto"

	"github.com/golang/protobuf/proto"
	"github.com/murphybytes/ucp/common"
)

type respondent interface {
	getMessage() (msg proto.Message, e error)
	respondMessage(msg proto.Message) (e error)
	setClientEncryptionKey(key crypto.PublicKey)
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
