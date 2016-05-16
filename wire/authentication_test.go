package wire

import (
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestAuthenticationRequest(t *testing.T) {
	request := &AuthenticationRequest{
		UserName:   "foo",
		MethodName: AuthenticationMethodPublicKey,
		PublicKey:  "DEADBEEF",
	}

	out, err := proto.Marshal(request)
	if err != nil {
		t.Error("Did not expect", err.Error())
	}

	response := &AuthenticationRequest{}
	err = proto.Unmarshal(out, response)
	if err != nil {
		t.Error("Did not expect error", err.Error())
	}

	if response.PublicKey != "DEADBEEF" {
		t.Error("What went in wasn't what came out.")
	}

}
