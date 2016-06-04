package wire

// ReadPacket buffer sent from server, contains bytes read from file
// and status of transfer
type ReadPacket struct {
	NextInitializationVector []byte
	Buffer                   []byte
	Status                   ResponseCode
	StatusText               string
}
