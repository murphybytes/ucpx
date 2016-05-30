package wire

import (
	"crypto/rsa"
)

// AuthenticationRequest initial request from client
type AuthenticationRequest struct {
	UserName                      string
	RequestedAuthenticationMethod string
	PublicKey                     rsa.PublicKey
}

// AutenticationResponse response to initial request
type AutenticationResponse struct {
	UserName                    string
	AllowedAuthenticationMethod string
	Status                      ResponseCode
	StatusText                  string
	PublicKey                   rsa.PublicKey
}
