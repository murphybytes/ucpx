package common

import (
	"errors"
	"flag"
	"os"
)

// Flags - persisted command line arguments
type Flags struct {
	// if true app runs as server
	IsServer bool
	// If true print help and exit
	Help bool
	// User name on host
	User string
}

// NewFlags returns a pointer to Flags which contains command line variables
func NewFlags() (flags *Flags, e error) {
	flags = &Flags{}
	flag.BoolVar(&flags.IsServer, "server", false, "Server mode. If set the application will listen for incoming client requests")
	// client options
	flag.StringVar(&flags.User, "user", "", "Client mode.  User on the target system.")
	flag.BoolVar(&flags.Help, "help", false, "Prints Usage")
	flag.Parse()

	if flags.IsServer {
		e = validateServerFlags(flags)
	} else {
		e = validateClientFlags(flags)
	}

	if flags.Help {
		flag.PrintDefaults()
		os.Exit(2)
	}
	return flags, e
}

func validateClientFlags(flags *Flags) (e error) {
	if flags.User == "" {
		e = errors.New("-user is required in client mode")
	}
	return e
}

func validateServerFlags(flags *Flags) (e error) {
	return e
}
