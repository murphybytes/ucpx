package client

import (
	"fmt"
	"io"
	"log"

	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/wire"
)

// Client contains all the logic for ucp client.
type Client struct {
	flags *common.Flags
}

// New create new Client
func New(flags *common.Flags) common.Application {
	return &Client{
		flags: flags,
	}
}

// Run the client application
func (c *Client) Run() (e error) {
	var reader io.ReadCloser
	var writer io.WriteCloser

	if reader, e = newReader(c.flags); e != nil {
		log.Fatal(e.Error())
	}
	defer reader.Close()

	if writer, e = newWriter(c.flags); e != nil {
		log.Fatal(e.Error())
	}
	defer writer.Close()

	fmt.Println("opened reader and writer")
	for {
		var read int
		readBuffer := make([]byte, wire.DataBufferSize)

		read, e = reader.Read(readBuffer)

		fmt.Printf("Read %d\n", read)

		if e == io.EOF {
			break
		}

		if e != nil {
			log.Fatal(e.Error())
		}

		fmt.Println("Ready to writexxx ")
		var written int
		if written, e = writer.Write(readBuffer[:read]); e != nil {
			fmt.Println("Got here")
			log.Fatal(e.Error())
		}

		fmt.Printf("client write successful wrote %d\n", written)

	}

	return
}
