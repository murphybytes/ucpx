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
	var logger common.Logger
	if logger, e = common.NewLogger(s.flags); e != nil {
		return
	}

	var listener net.Listener
	connectString := getServerString(s.flags)
	logger.LogInfo("Connecting to ", connectString)
	if listener, e = udt.Listen(connectString); e != nil {
		return
	}

	defer listener.Close()

	for connectionCount := int64(1); ; connectionCount++ {
		var conn net.Conn
		conn, e = listener.Accept()

		if e == nil {

			go handleConnection(context{
				flags:  s.flags,
				conn:   conn,
				logger: logger,
				connID: connectionCount,
			})

		} else {
			logger.LogError(e.Error())
			break
		}
	}

	return
}

func handleConnection(ctx context) {
	defer ctx.conn.Close()
	ctx.logger.LogInfo("Connection from ", ctx.conn.RemoteAddr())

	var e error
	var client respondent
	if client, e = newClient(&ctx); e != nil {
		ctx.logger.LogError("Client creation failed -", e.Error())
	}

	if e = authenticate(client, &ctx); e != nil {
		ctx.logger.LogError("Authentication failed -", e.Error())
		return
	}

}
