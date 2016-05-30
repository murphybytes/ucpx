package client

import (
	"crypto"
	"fmt"
	"net"

	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/udt.go/udt"
)

type context struct {
	fileInfo  *fileInfo
	flags     *common.Flags
	logger    common.Logger
	server    requester
	publicKey crypto.PublicKey
}

func getReadContext(flags *common.Flags) (c *context, e error) {

	return getContext(flags.From, flags, true)
}

func getWriteContext(flags *common.Flags) (c *context, e error) {
	return getContext(flags.To, flags, false)
}

func getContext(filespec string, flags *common.Flags, read bool) (c *context, e error) {
	var fi *fileInfo
	fi, e = newFileInfo(filespec, read)

	if e != nil {
		return
	}

	var logger common.Logger
	if logger, e = common.NewLogger(flags); e != nil {
		return
	}

	ctx := &context{
		fileInfo: fi,
		flags:    flags,
		logger:   logger,
	}

	if !fi.local {
		var connectString string
		connectString, e = fi.getConnectString()
		logger.LogInfo("Client connecting to ", connectString)
		if e != nil {
			return
		}
		var conn net.Conn
		conn, e = udt.Dial(connectString)
		if e != nil {
			return
		}

		if ctx.server, e = newServer(conn, ctx); e != nil {
			return
		}

		e = auth(ctx, fmt.Scanln)

	}

	return
}

func (c *context) close() {
	if c.server != nil {
		c.server.close()
	}
}
