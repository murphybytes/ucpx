package main

import (
	"fmt"
	"os"

	"github.com/murphybytes/ucp/client"
	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/server"
)

func main() {
	var app common.Application
	flags := common.NewFlags()
	if flags.IsServer {
		app = server.New()
	} else {
		app = client.New()
	}
	err := app.Run(flags)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)

}
