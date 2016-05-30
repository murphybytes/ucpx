package common

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	missingSourceMessage  = "-from is required in client mode."
	missingTargetMessage  = "-to is required in client mode."
	invalidLogLevel       = "-verbosity argument is not valid, must be one of INFO WARN ERROR"
	missingPublicKeyPath  = "-public-key-path is required"
	missingPrivateKeyPath = "-private-key-path is required"

	logInfo  = "INFO"
	logWarn  = "WARN"
	logError = "ERROR"
	// DefaultPort used to connect unless changed on command line
	DefaultPort = 9191
	// KeySize default RSA key size
	KeySize = 4096
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
	// Port of ucp server
	Port int
	// Interface the server uses
	Host string
	// Log level INFO, WARN, ERROR
	LogLevel string
	// Path to public crypto key
	PublicKeyPath string
	//  Path to private key
	PrivateKeyPath string
	// Generate public private keys and exit
	GenerateKeys bool
}

// NewFlags returns a pointer to Flags which contains command line variables
func NewFlags() (flags *Flags) {
	var e error
	flags = &Flags{}
	flag.BoolVar(&flags.IsServer, "server", false, "Server mode. If set the application will listen for incoming client requests")
	// client options
	flag.StringVar(&flags.From, "from", "", "Client mode file to copy from.  [[user]@[host]:]filepath")
	flag.StringVar(&flags.To, "to", "", "Client mode file to copy to. [[user]@[host]:]filepath")
	flag.IntVar(&flags.Port, "port", DefaultPort, "Server Mode. The port that the ucp server listens on")
	flag.StringVar(&flags.Host, "host", "127.0.0.1", "Server Mode. The host or interface the server listens on")
	flag.StringVar(&flags.LogLevel, "verbosity", logWarn, "Log level. INFO|WARN|ERROR")
	flag.StringVar(&flags.PrivateKeyPath, "private-key-path", getDefaultKeyPath("ucp.pem"), "Path to private key")
	flag.StringVar(&flags.PublicKeyPath, "public-key-path", getDefaultKeyPath("key.pub"), "Path to public key")
	flag.BoolVar(&flags.GenerateKeys, "generate-keys", false, "Generate key pair and exit")
	flag.BoolVar(&flags.Help, "help", false, "Prints Usage")
	flag.Parse()

	if e = validateFlags(flags); e != nil {
		fmt.Println("Missing or invalid command line arguments -", e.Error())
		fmt.Println()
		flag.PrintDefaults()
		os.Exit(1)
	}

	if flags.Help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if flags.GenerateKeys {
		if e = ucpKeyGenerate(flags.PrivateKeyPath, flags.PublicKeyPath); e == nil {
			fmt.Println("Key generation successful")
			fmt.Println("Public key ->", flags.PublicKeyPath)
			fmt.Println("Private key ->", flags.PrivateKeyPath)
			os.Exit(0)
		} else {
			fmt.Println("Key generation failed -", e.Error())
			os.Exit(1)
		}
	}

	return flags
}

func getDefaultKeyPath(keyname string) (path string) {
	if homeDir := os.Getenv("HOME"); homeDir != "" {
		path = fmt.Sprintf("%s/.ucp/%s", homeDir, keyname)
	}
	return

}

func validateClientFlags(flags *Flags) (e error) {
	if flags.From == "" {
		e = errors.New(missingSourceMessage)
		return
	}

	if flags.To == "" {
		e = errors.New(missingTargetMessage)
		return
	}

	return
}

func validateKeygenFlags(flags *Flags) error {
	if flags.PrivateKeyPath == "" {
		return errors.New(missingPrivateKeyPath)
	}

	if flags.PublicKeyPath == "" {
		return errors.New(missingPublicKeyPath)
	}

	return nil
}

func validateFlags(flags *Flags) (e error) {
	flags.LogLevel = strings.ToUpper(flags.LogLevel)
	if !(flags.LogLevel == logInfo || flags.LogLevel == logWarn || flags.LogLevel == logError) {
		e = errors.New(invalidLogLevel)
		return
	}

	if flags.GenerateKeys {
		return validateKeygenFlags(flags)
	}

	if flags.IsServer {
		return validateServerFlags(flags)
	}

	return validateClientFlags(flags)
}

func validateServerFlags(flags *Flags) (e error) {
	return
}
