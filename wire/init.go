package wire

type TransferType int

const (
	ClientReading TransferType = iota
	ClientWriting
)

type FileTransferRequest struct {
	UserName string
	FilePath string
	Transfer TransferType
}

type FileTransferResponse struct {
	Status               ResponseCode
	StatusText           string
	AESKey               []byte
	InitializationVector []byte
}
