package client

import (
	"log"

	"github.com/murphybytes/ucp/common"
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
	var reader common.Reader
	var writer common.Writer

	if reader, e = newReader(c.flags); e != nil {
		log.Fatal(e.Error())
	}
	defer reader.Close()

	if writer, e = newWriter(c.flags); e != nil {
		log.Fatal(e.Error())
	}
	defer writer.Close()

	// for {
	// 	var buffer []byte
	//
	// 	if buffer, e = reader.Read(); e != nil {
	// 		if e == io.EOF {
	// 			break
	// 		} else {
	// 			log.Fatal(e.Error())
	// 		}
	// 	}
	//
	// 	if e = writer.Write(buffer); e != nil {
	// 		log.Fatal(e.Error())
	// 	}
	//
	// }

	return
}
