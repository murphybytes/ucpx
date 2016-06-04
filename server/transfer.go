package server

import (
	"crypto/cipher"
	"io"
)

type transferContext struct {
	fileTransferRequest  *wire.FileTransferRequest
	block                *cipher.Block
	initializationVector []byte
	conn                 io.ReadWriteCloser
}

func readFromClient(ctx *transferContext) (e error) {

	return
}

func writeToClient(ctx *transferContext) (e error) {
	return
}
