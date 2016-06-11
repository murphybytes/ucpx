package server

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"io"

	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/wire"
)

type transferContext struct {
	block                cipher.Block
	initializationVector []byte
	conn                 io.ReadWriteCloser
}

func readRemoteWriteLocal(ctx *transferContext, outfile io.Writer) (e error) {

	for {
		var read int
		encrypted := make([]byte, wire.ReadBufferSize)
		if read, e = ctx.conn.Read(encrypted); e != nil {
			return
		}

		decrypted := common.DecryptAES(ctx.block, ctx.initializationVector, encrypted[:read])

		decoderBuffer := bytes.NewBuffer(decrypted)
		decoder := gob.NewDecoder(decoderBuffer)

		clientRead := &wire.ClientRead{}
		if e = decoder.Decode(clientRead); e != nil {
			return
		}

		if clientRead.Status == wire.EOF {
			e = io.EOF
			return
		}

		newIV := make([]byte, common.IVBlockSize)
		rand.Read(newIV)

		response := wire.ClientReadResponse{
			NextInitializationVector: newIV,
		}

		var err error
		if _, e = outfile.Write(clientRead.Buffer); e != nil {
			// Tell client to stop sending and disconnect
			response.Status = wire.Error
			response.StatusText = e.Error()
			err = e
		} else {
			response.Status = wire.OK
			response.StatusText = "OK"
		}

		var encoderBuffer bytes.Buffer
		encoder := gob.NewEncoder(&encoderBuffer)

		if e = encoder.Encode(response); e != nil {
			return
		}

		encrypted = common.EncryptAES(ctx.block, ctx.initializationVector, encoderBuffer.Bytes())

		if _, e = ctx.conn.Write(encrypted); e != nil || err != nil {
			if e == nil {
				e = err
			}
			return
		}

		// client will user newIV for message they send back to us
		ctx.initializationVector = newIV

	}

}

func readLocalWriteRemote(ctx *transferContext, infile io.Reader) (e error) {

	for {

		var read int
		encrypted := make([]byte, wire.ReadBufferSize)

		if read, e = ctx.conn.Read(encrypted); e != nil {
			return
		}

		decrypted := common.DecryptAES(ctx.block, ctx.initializationVector, encrypted[:read])

		decodeBuffer := bytes.NewBuffer(decrypted)
		decoder := gob.NewDecoder(decodeBuffer)
		clientDataRequest := &wire.ClientDataRequest{}
		if e = decoder.Decode(clientDataRequest); e != nil {
			return
		}

		if clientDataRequest.Status != wire.More {
			return errors.New(clientDataRequest.StatusText)
		}

		newIV := make([]byte, common.IVBlockSize)
		rand.Read(newIV)

		data := make([]byte, wire.DataBufferSize)

		if read, e = infile.Read(data); e != nil {
			// send message to client to terminate connection
			var empty []byte

			if e == io.EOF {
				return sendClientDataResponse(ctx, newIV, empty, wire.EOF, "End of data")
			}

			sendClientDataResponse(ctx, newIV, empty, wire.Error, e.Error())
			return

		}

		if e = sendClientDataResponse(ctx, newIV, data[:read], wire.OK, "OK"); e != nil {
			return
		}

	}

}

func sendClientDataResponse(ctx *transferContext, iv []byte, data []byte, status wire.ResponseCode, statusText string) (e error) {
	// note we send the next iv to client who will use it to encrypt the next message sent
	// to us.  We use ctx.initializationVector to ecrypt this message
	response := wire.ClientDataResponse{
		NextInitializationVector: iv,
		Data:       data,
		Status:     status,
		StatusText: statusText,
	}

	var encoderBuffer bytes.Buffer
	encoder := gob.NewEncoder(&encoderBuffer)
	if e = encoder.Encode(response); e != nil {
		return
	}

	encrypted := common.EncryptAES(ctx.block, ctx.initializationVector, encoderBuffer.Bytes())

	if _, e = ctx.conn.Write(encrypted); e != nil {
		return
	}
	// set iv that client will use to encrypt the next message it sends
	ctx.initializationVector = iv

	return
}
