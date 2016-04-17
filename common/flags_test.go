package common

import "testing"

func TestClientValidation(t *testing.T) {
	flags := Flags{}
	err := validateClientFlags(&flags)
	if err == nil {
		t.Error("Expecting client validation error")
	}
}
