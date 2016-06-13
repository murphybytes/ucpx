package server

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"io"
	mathrand "math/rand"
	"testing"
	"time"

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

type stringable interface {
	toString() string
}

type waitable interface {
	wait()
}

// mock read closer to read from server
type mockServerReadFile struct {
	fileBytes []byte
	read      int
}

func (m *mockServerReadFile) toString() string {
	return string(m.fileBytes)
}

func newMockServerReadFile(testFileSize int) io.ReadCloser {
	var mock mockServerReadFile
	//make a fake file full of random crap
	mock.fileBytes = make([]byte, testFileSize)
	rand.Read(mock.fileBytes)

	return &mock
}

func (m *mockServerReadFile) Read(b []byte) (n int, e error) {
	if m.read == len(m.fileBytes) {
		return 0, io.EOF
	}

	n = copy(b, m.fileBytes[m.read:])

	m.read += n
	return
}

func (m *mockServerReadFile) Close() (e error) {
	return
}

// represents remote client io.ReadWriteCloser
type mockClientRecipient struct {
	readNum        int
	sentFromServer []byte
	fauxNetwork    chan []byte
	waiter         chan int
}

func (m *mockClientRecipient) toString() string {
	return string(m.sentFromServer)
}

func (m *mockClientRecipient) wait() {
	<-m.waiter
}

func getMockClientReaderContext() (ctx *transferContext) {

	iv := make([]byte, common.IVBlockSize)

	rand.Read(iv)

	block, _ := common.NewCipherBlock()

	mockClient := &mockClientRecipient{
		fauxNetwork: make(chan []byte),
		waiter:      make(chan int),
	}

	ctx = &transferContext{
		block:                block,
		initializationVector: iv,
		conn:                 mockClient,
	}

	go func() {
		clientiv := iv
		clientblock := block
		defer func() {
			mockClient.waiter <- 1
		}()

		for {

			clientDataRequest := &wire.ClientDataRequest{
				Status:     wire.More,
				StatusText: "More",
			}

			var encodeBuffer bytes.Buffer
			encoder := gob.NewEncoder(&encodeBuffer)
			if err := encoder.Encode(clientDataRequest); err != nil {
				return
			}

			encrypted := common.EncryptAES(clientblock, clientiv, encodeBuffer.Bytes())
			mockClient.fauxNetwork <- encrypted
			responseBuffer := <-mockClient.fauxNetwork

			decrypted := common.DecryptAES(clientblock, clientiv, responseBuffer)

			decodeBuffer := bytes.NewBuffer(decrypted)
			decoder := gob.NewDecoder(decodeBuffer)
			var response wire.ClientDataResponse
			if err := decoder.Decode(&response); err != nil {
				return
			}

			if response.Status == wire.EOF {
				return
			}

			if response.Status != wire.OK {
				return
			}

			encrypted = <-mockClient.fauxNetwork
			decrypted = common.DecryptAES(clientblock, clientiv, encrypted)
			mockClient.sentFromServer = append(mockClient.sentFromServer, decrypted...)

			clientiv = response.NextInitializationVector

		}

	}()

	return

}

func (m *mockClientRecipient) Read(b []byte) (n int, e error) {
	response := <-m.fauxNetwork
	n = copy(b, response)
	return
}

func (m *mockClientRecipient) Write(b []byte) (n int, e error) {
	n = len(b)
	m.fauxNetwork <- b
	return
}

func (m *mockClientRecipient) Close() (e error) {
	return
}

func TestReadLocalWriteRemote(t *testing.T) {
	mathrand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		filesize := mathrand.Intn(0x100000)

		file := newMockServerReadFile(filesize)
		ctx := getMockClientReaderContext()
		defer ctx.conn.Close()
		if err := readLocalWriteRemote(ctx, file); err != nil {
			t.Fatal("Unexpected error ", err.Error())
		}

		ctx.conn.(waitable).wait()
		client := ctx.conn.(stringable)
		serverfile := file.(stringable)

		if client.toString() != serverfile.toString() {
			t.Fatal("Contents should match got ", len(client.toString()), " expected ", len(serverfile.toString()))
		}

	}

}
