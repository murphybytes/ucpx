package client

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/wire"
	"github.com/murphybytes/udt.go/udt"
)

type context struct {
	fileInfo             *fileInfo
	flags                *common.Flags
	logger               common.Logger
	server               requester
	file                 *os.File
	publicKey            crypto.PublicKey
	aesKey               cipher.Block
	initializationVector []byte
}

func getReadContext(flags *common.Flags) (c *context, e error) {

	return getContext(flags.From, flags, true)
}

func getWriteContext(flags *common.Flags) (c *context, e error) {
	return getContext(flags.To, flags, false)
}

func getContext(filespec string, flags *common.Flags, read bool) (c *context, e error) {
	var fi *fileInfo
	fi, e = newFileInfo(filespec, read)

	if e != nil {
		return
	}

	var logger common.Logger
	if logger, e = common.NewLogger(flags); e != nil {
		return
	}

	ctx := &context{
		fileInfo: fi,
		flags:    flags,
		logger:   logger,
	}

	if !fi.local {
		// remote context read or write encrypted bytes to a socket
		var connectString string
		connectString, e = fi.getConnectString()
		logger.LogInfo("Client connecting to ", connectString)
		if e != nil {
			return
		}
		var conn net.Conn
		conn, e = udt.Dial(connectString)
		if e != nil {
			return
		}

		if ctx.server, e = newServer(conn, ctx); e != nil {
			return
		}

		if e = auth(ctx, fmt.Scanln); e != nil {
			return
		}

		e = initTransfer(ctx)

	} else {
		// local context read or write to a file
		if ctx.fileInfo.read {
			if ctx.file, e = os.Open(fi.path); e != nil {
				return
			}
		} else {
			if ctx.file, e = os.Create(fi.path); e != nil {
				return
			}

		}
	}

	return
}

func initTransfer(ctx *context) (e error) {

	var direction wire.TransferType
	if ctx.fileInfo.read {
		direction = wire.ClientReading
	} else {
		direction = wire.ClientWriting
	}

	txfrRequest := wire.FileTransferRequest{
		UserName: ctx.fileInfo.user,
		FilePath: ctx.fileInfo.path,
		Transfer: direction,
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if e = encoder.Encode(txfrRequest); e != nil {
		return
	}

	var response []byte
	if response, e = ctx.server.get(buffer.Bytes()); e != nil {
		return
	}

	decoderBuffer := bytes.NewBuffer(response)
	decoder := gob.NewDecoder(decoderBuffer)
	var txfrResponse wire.FileTransferResponse
	if e = decoder.Decode(&txfrResponse); e != nil {
		return
	}

	if txfrResponse.Status != wire.OK {
		return errors.New(txfrResponse.StatusText)
	}

	if ctx.aesKey, e = aes.NewCipher(txfrResponse.AESKey); e != nil {
		return
	}

	ctx.initializationVector = txfrResponse.InitializationVector

	return

}

func (c *context) getIO() io.ReadWriteCloser {
	if c.server != nil {
		return c.server
	}

	return c.file
}

func (c *context) Read(p []byte) (n int, e error) {
	reader := c.getIO()
	return reader.Read(p)

}

func (c *context) Write(p []byte) (n int, e error) {
	writer := c.getIO()
	return writer.Write(p)
}

//func (c *context)
func (c *context) Close() error {
	closer := c.getIO()
	return closer.Close()

}
