package server

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"io"
	"testing"

	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/wire"
)

type mockClientConn struct {
	sent        []byte
	readNumber  int
	writeNumber int
	block       cipher.Block
	iv          []byte
}

type mockServerFile struct {
	writeNumber int
	received    []byte
}

func (m *mockClientConn) Read(b []byte) (n int, e error) {

	response := wire.ClientRead{
		Status:     wire.More,
		StatusText: "More data is available",
	}

	if m.readNumber == 0 {
		response.Buffer = m.sent[:1024]
		m.readNumber++
	} else if m.readNumber == 1 {
		response.Buffer = m.sent[1024:]
		m.readNumber++
	} else if m.readNumber == 2 {
		response.Status = wire.EOF
		response.StatusText = "End of data"
		m.readNumber++
	} else {
		return 0, errors.New("Shouldn't have three reads something is wrong")
	}

	var encodeBuffer bytes.Buffer
	encoder := gob.NewEncoder(&encodeBuffer)
	if e = encoder.Encode(response); e != nil {
		return
	}

	encrypted := common.EncryptAES(m.block, m.iv, encodeBuffer.Bytes())
	n = copy(b, encrypted)

	return
}

func (m *mockClientConn) Write(b []byte) (n int, e error) {
	n = len(b)
	decrypted := common.DecryptAES(m.block, m.iv, b)
	decodeBuffer := bytes.NewBuffer(decrypted)
	decoder := gob.NewDecoder(decodeBuffer)
	var serverResponse wire.ClientReadResponse
	if e = decoder.Decode(&serverResponse); e != nil {
		return
	}

	m.iv = serverResponse.NextInitializationVector

	return
}

func (m *mockClientConn) Close() (e error) {
	return
}

func (m *mockServerFile) Write(b []byte) (n int, e error) {
	m.received = append(m.received, b...)
	return
}

func getReadRemoteWriteLocalMock(sendBuffer []byte) (ctx *transferContext, f *mockServerFile) {

	iv := make([]byte, common.IVBlockSize)

	rand.Read(iv)

	block, _ := common.NewCipherBlock()

	ctx = &transferContext{
		block:                block,
		initializationVector: iv,
		conn: &mockClientConn{
			sent:  sendBuffer,
			block: block,
			iv:    iv,
		},
	}

	f = &mockServerFile{
		received: []byte{},
	}

	return
}

func TestReadRemoteWriteLocal(t *testing.T) {
	sendBuffer := make([]byte, 2048)
	rand.Read(sendBuffer)
	ctx, file := getReadRemoteWriteLocalMock(sendBuffer)
	if e := readRemoteWriteLocal(ctx, file); e != io.EOF {
		t.Fatal("readRemoteWriteLocal should be EOF - ", e.Error())
	}

	if string(file.received) != string(sendBuffer) {
		t.Fatal("We didn't get the buffer that was sent")
	}
}
