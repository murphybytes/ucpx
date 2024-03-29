package wire

const (
	// AuthenticationMethodPublicKey public key method
	AuthenticationMethodPublicKey = "PUBLIC_KEY"
	//AuthenticationMethodPassword password method
	AuthenticationMethodPassword = "PASSWORD"

	// TxferBufferSize is the size of a preallocated buffer we use for network
	// file reads and writes
	TxferBufferSize = 0x10000
	// DataBufferSize size of data packet for ClientRead and ClientDataResponse
	DataBufferSize = 0x10000
	// ReadBufferSize size of buffer read from network connection
	ReadBufferSize = 0x2800 + DataBufferSize
)

// ResponseCode codes to communicate status of transactions
type ResponseCode int

const (
	// OK indicates successful transaction
	OK ResponseCode = iota
	// Error indicates failure check status text for details
	Error
	// MethodNotAllowed invalid command was recieved
	MethodNotAllowed
	// EOF end of file, no more bytes available to be read
	EOF
	// More more data to read from client
	More
)
