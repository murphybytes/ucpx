package client

import (
	"github.com/murphybytes/ucp/common"
)

type reader struct {
	ctx *context
}

func newReader(filespec string) (r common.Reader, e error) {
	var ctx *context
	ctx, e = getContext(filespec)
	if e != nil {
		return
	}
	r = &reader{
		ctx: ctx,
	}

	return
}

func (r *reader) Read() (buff []byte, e error) {
	return
}

func (r *reader) Close() {

}
