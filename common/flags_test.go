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

func TestLoggingArg(t *testing.T) {
	flags := &Flags{
		IsServer: true,
		LogLevel: "soemthing",
	}
	err := validateFlags(flags)
	if err == nil {
		t.Error("Expected error because LogLevel is invalid")
	}

	flags.LogLevel = logInfo
	err = validateFlags(flags)
	if err != nil {
		t.Error("We don't expect error when log level is valid")
	}
}
