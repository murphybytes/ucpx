package wire

const (
	// AuthenticationMethodPublicKey public key method
	AuthenticationMethodPublicKey = "PUBLIC_KEY"
	//AuthenticationMethodPassword password method
	AuthenticationMethodPassword = "PASSWORD"

	// ReadBufferSize size of buffer read from network connection
	ReadBufferSize = 0x2800
)

type ResponseCode int

const (
	OK ResponseCode = iota
	Error
	MethodNotAllowed
)
