package main

import (
	"log"

	"github.com/murphybytes/ucp/client"
	"github.com/murphybytes/ucp/common"
	"github.com/murphybytes/ucp/server"
	"github.com/murphybytes/udt.go/udt"
)

func main() {
	udt.Startup()
	defer udt.Cleanup()
	flags := common.NewFlags()
	app := newApplication(flags)
	err := app.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

}

func newApplication(f *common.Flags) (app common.Application) {
	if f.IsServer {
		app = server.New(f)
	} else {
		app = client.New(f)
	}
	return app
}
