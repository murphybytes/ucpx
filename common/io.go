package common

// Reader abstracts reading a file which
// can be local or remote
type Reader interface {
	Read() ([]byte, error)
	Close()
}

// Writer abstracts writing a file which
// can be local or remote
type Writer interface {
	Write(buffer []byte) error
	Close()
}

// Authenticator handles authenticating operation
type Authenticator interface {
	Authenticate() (bool, error)
}
