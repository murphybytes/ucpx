package client

import (
	"github.com/murphybytes/ucp/common"
)

// Client contains all the logic for ucp client.
type Client struct {
}

// New create new Client
func New() common.Application {
	return &Client{}
}

// Run the client application
func (c *Client) Run(flags *common.Flags) (e error) {
	return
}
