package common

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

const (
	missingSourceMessage = "-from is required in client mode."
	missingTargetMessage = "-to is required in client mode."
)

// Flags - persisted command line arguments
type Flags struct {
	// if true app runs as server
	IsServer bool
	// If true print help and exit
	Help bool
	// From name of file to copy from
	From string
	// To name of file to copy to
	To string
}

// NewFlags returns a pointer to Flags which contains command line variables
func NewFlags() (flags *Flags) {
	var e error
	flags = &Flags{}
	flag.BoolVar(&flags.IsServer, "server", false, "Server mode. If set the application will listen for incoming client requests")
	// client options
	flag.StringVar(&flags.From, "from", "", "Client mode file to copy from.  [[user]@[host]:]filepath")
	flag.StringVar(&flags.To, "to", "", "Client mode file to copy to. [[user]@[host]:]filepath")
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

	if e != nil {
		fmt.Println(e)
		flag.PrintDefaults()
		os.Exit(2)
	}

	return flags
}

func validateClientFlags(flags *Flags) (e error) {
	if flags.From == "" {
		e = errors.New(missingSourceMessage)
		return e
	}

	if flags.To == "" {
		e = errors.New(missingTargetMessage)
		return e
	}

	return e
}

func validateServerFlags(flags *Flags) (e error) {
	return e
}
