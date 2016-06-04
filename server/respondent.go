package server

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/gob"
	"fmt"
	"math/big"
	"os/user"

	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/wire"
)

type respondent interface {
	initializeSecureChannel() (e error)
	initializeTransfer() (e error)
	getMessage() ([]byte, error)
	sendMessage([]byte) (e error)
}

type client struct {
	clientKey    *rsa.PublicKey
	serverKey    *rsa.PrivateKey
	context      *context
	transferInfo *wire.FileTransferRequest
	aesKey       cipher.Block
	startingIV   []byte
}

func newClient(ctx *context) (r respondent, e error) {

	r = &client{

		context: ctx,
	}

	return

}

// first message from client is unencrypted and contains their public key
func (c *client) initializeSecureChannel() (e error) {
	c.context.logger.LogInfo("Beginning public key exchange with client")
	networkReadBuff := make([]byte, wire.ReadBufferSize)
	var readBytes int

	if readBytes, e = c.context.conn.Read(networkReadBuff); e != nil {
		return
	}

	encoderBuffer := bytes.NewBuffer(networkReadBuff[:readBytes])
	decoder := gob.NewDecoder(encoderBuffer)

	authRequest := wire.AuthenticationRequest{
		PublicKey: rsa.PublicKey{
			N: &big.Int{},
		},
	}

	if e = decoder.Decode(&authRequest); e != nil {
		return
	}

	c.context.logger.LogInfo("Successfully received authRequest for ", authRequest.UserName)

	if c.serverKey, e = getUserPrivateKey(authRequest.UserName); e != nil {
		return
	}

	// TODO: check authorization here

	c.clientKey = &authRequest.PublicKey

	// we now have clients public key, so send server public key to client

	authResponse := wire.AutenticationResponse{
		UserName:                    authRequest.UserName,
		PublicKey:                   c.serverKey.PublicKey,
		AllowedAuthenticationMethod: authRequest.RequestedAuthenticationMethod,
		Status:     wire.OK,
		StatusText: "OK",
	}

	encoderBuffer.Reset()
	encoder := gob.NewEncoder(encoderBuffer)

	if e = encoder.Encode(authResponse); e != nil {
		return
	}

	if _, e = c.context.conn.Write(encoderBuffer.Bytes()); e != nil {
		return
	}

	c.context.logger.LogInfo("Sent our public key to client")

	return
}

/////////////////////////////////////////////////
// exchange keys
func (c *client) initializeTransfer() (e error) {
	c.context.logger.LogInfo("Preparing for file transfer")

	var clientMsg []byte
	if clientMsg, e = c.getMessage(); e != nil {
		return
	}

	decoderBuffer := bytes.NewBuffer(clientMsg)
	decoder := gob.NewDecoder(decoderBuffer)

	c.transferInfo = &wire.FileTransferRequest{}

	if e = decoder.Decode(c.transferInfo); e != nil {
		return
	}

	c.context.logger.LogInfo("Recieved trasfer message preparing AES key")
	// generate random key and initialization vector for aes-256
	keylen := 32
	keybuff := make([]byte, keylen)
	if _, e = rand.Read(keybuff); e != nil {
		return
	}

	if c.aesKey, e = aes.NewCipher(keybuff); e != nil {
		return
	}

	c.startingIV = make([]byte, aes.BlockSize)
	if _, e = rand.Read(c.startingIV); e != nil {
		return
	}

	fileTxfrResponse := wire.FileTransferResponse{
		Status:               wire.OK,
		StatusText:           "OK",
		AESKey:               keybuff,
		InitializationVector: c.startingIV,
	}

	var encodeBuff bytes.Buffer
	encoder := gob.NewEncoder(&encodeBuff)
	if e = encoder.Encode(fileTxfrResponse); e != nil {
		return
	}

	if e = c.sendMessage(encodeBuff.Bytes()); e != nil {
		return
	}

	c.context.logger.LogInfo("Sent AES key to client")

	return

}

func (c *client) getMessage() (msg []byte, e error) {
	buffer := make([]byte, wire.ReadBufferSize)
	var read int
	if read, e = c.context.conn.Read(buffer); e != nil {
		return
	}

	if msg, e = common.DecryptOAEP(c.serverKey, buffer[:read]); e != nil {
		return
	}

	return msg, nil

}

func (c *client) sendMessage(msg []byte) (e error) {

	var encrypted []byte
	if encrypted, e = common.EncryptOAEP(c.clientKey, msg); e != nil {
		return
	}

	_, e = c.context.conn.Write(encrypted)

	return
}

func getUserPrivateKey(userName string) (key *rsa.PrivateKey, e error) {
	var u *user.User
	if u, e = user.Lookup(userName); e != nil {
		return
	}

	privateKeyPath := fmt.Sprint(u.HomeDir, "/.ucp/ucp.pem")

	key, e = common.GetPrivateKey(privateKeyPath)

	return

}
