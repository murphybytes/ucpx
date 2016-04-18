package common

import "testing"

func TestClientValidation(t *testing.T) {
	flags := Flags{}
	err := validateClientFlags(&flags)
	if err == nil {
		t.Error("Expecting client validation error")
	} else if err.Error() != missingSourceMessage {
		t.Error("Expected ", missingSourceMessage, " got ", err)

	}

	flags.From = "foo@bar:/zip"
	err = validateClientFlags(&flags)
	if err == nil {
		t.Error("Expecting client validation error, to is missing")
	} else if err.Error() != missingTargetMessage {
		t.Error("Expected ", missingTargetMessage, " got ", err)
	}
}
