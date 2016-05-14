package server

import (
	"fmt"
	"net"

	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/udt.go/udt"
)

// Server contains all the logic involved with handling incoming client
// connections
type Server struct {
	flags *common.Flags
}

// New creates a Server
func New(flags *common.Flags) common.Application {
	return &Server{
		flags: flags,
	}
}

func getServerString(flags *common.Flags) string {
	return fmt.Sprintf("%s:%d", flags.Host, flags.Port)
}

// Run the application as a server
func (s *Server) Run() (e error) {
	var listener net.Listener
	if listener, e = udt.Listen(getServerString(s.flags)); e == nil {
		var logger common.Logger

		if logger, e = common.NewLogger(s.flags); e == nil {
			for count := 1; ; count++ {
				var conn net.Conn
				conn, e = listener.Accept()
				if e == nil {
					defer conn.Close()
					go handleConnection(context{
						flags:  s.flags,
						conn:   conn,
						logger: logger,
						connID: count,
					})
				} else {
					logger.LogError(e.Error())
				}
			}
		}

	}
	return
}

func handleConnection(ctx context) {
	ctx.logger.LogInfo("Connection", ctx.connID, "starting...")

}
