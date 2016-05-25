# ucp
Super fast encrypted network file copy that uses [UDT](http://udt.sourceforge.net/), A UDP based network protocol optimized for bulk data transfers. 

## Usage

Generate key pair that will be used to encrypt communication. For example, this command will generate a key pair in the default location.  
```
ucp -generate-keys
```

### Command Line Options

```
  jam [master] $ ucp --help
  -from string
        Client mode file to copy from.  [[user]@[host]:]filepath
  -generate-keys
        Generate key pair and exit
  -help
        Prints Usage
  -host string
        Server Mode. The host or interface the server listens on (default "localhost")
  -port int
        Server Mode. The port that the ucp server listens on (default 9191)
  -private-key-path string
        Path to private key (default "/Users/jam/.ucp/private.pem")
  -public-key-path string
        Path to public key (default "/Users/jam/.ucp/public.pem")
  -server
        Server mode. If set the application will listen for incoming client requests
  -to string
        Client mode file to copy to. [[user]@[host]:]filepath
  -verbosity string
        Log level. INFO|WARN|ERROR (default "WARN")
```

## Additional Documentation

Interaction between ucp client and server is documented in doc/wire.md.  
