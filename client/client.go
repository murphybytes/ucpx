package client

import (
	"github.com/murphybytes/ucp/common"
)

// Client contains all the logic for ucp client.
type Client struct {
}

// New create new Client
func New(flags *common.Flags) common.Application {
	return &Client{}
}

// Run the client application
func (c *Client) Run() (e error) {
	return
}
