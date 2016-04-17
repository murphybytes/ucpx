package server

import (
	"github.com/murphybytes/ucp/common"
)

// Server contains all the logic involved with handling incoming client
// connections
type Server struct {
}

// New creates a Server
func New() common.Application {
	return &Server{}
}

// Run the application as a server, flags contains command line
// arguments
func (s *Server) Run(flags *common.Flags) (e error) {
	return
}
