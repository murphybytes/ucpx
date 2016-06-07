package wire

// ReadPacket buffer sent from server, contains bytes read from file
// and status of transfer
type ReadPacket struct {
	NextInitializationVector []byte
	Buffer                   []byte
	Status                   ResponseCode
	StatusText               string
}

// ClientRead contains bytes read from client and sent to server.  If status is
// not EOF Buffer field contains bytes read from client file
type ClientRead struct {
	Buffer     []byte
	Status     ResponseCode
	StatusText string
}

// ClientReadResponse is returned to the client from the server.  The structure
// contains the initialization vector used to encrypt the next ClientRead
type ClientReadResponse struct {
	NextInitializationVector []byte
	Status                   ResponseCode
	StatusText               string
}
