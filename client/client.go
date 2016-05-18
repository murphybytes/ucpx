package client

import (
	"io"
	"log"
	"net"

	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/udt.go/udt"
)

// Client contains all the logic for ucp client.
type Client struct {
	flags *common.Flags
}

type context struct {
	conn     net.Conn
	fileInfo *fileInfo
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

	if reader, e = newReader(c.flags.From); e != nil {
		log.Fatal(e.Error())
	}
	defer reader.Close()

	if writer, e = newWriter(c.flags.To); e != nil {
		log.Fatal(e.Error())
	}
	defer writer.Close()

	for {
		var buffer []byte
		if buffer, e = reader.Read(); e != nil {
			if e == io.EOF {
				break
			} else {
				log.Fatal(e.Error())
			}
		}

		if e = writer.Write(buffer); e != nil {
			log.Fatal(e.Error())
		}
	}

	return
}

func getContext(filespec string) (c *context, e error) {
	var fi *fileInfo
	fi, e = newFileInfo(filespec)

	if e != nil {
		return
	}

	c = &context{
		fileInfo: fi,
	}

	if !fi.local {
		var connectString string
		connectString, e = fi.getConnectString()
		if e != nil {
			return
		}
		c.conn, e = udt.Dial(connectString)
		if e != nil {
			return
		}
		e = authenticate(c.conn, fi.user)
		if e != nil {
			c.conn.Close()
			return
		}
	}

	return
}
