package server

import (
	"net"

	"github.com/murphybytes/ucp/common"
)

type context struct {
	flags  *common.Flags
	conn   net.Conn
	logger common.Logger
	connID int
}

func newContext(flags *common.Flags, conn net.Conn) *context {
	return &context{
		flags: flags,
		conn:  conn,
	}
}
