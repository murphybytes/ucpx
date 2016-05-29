package client

import (
	"github.com/murphybytes/ucp/common"
)

type reader struct {
	ctx *context
}

func newReader(flags *common.Flags) (r common.Reader, e error) {
	var ctx *context
	ctx, e = getReadContext(flags)
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
	if r.ctx != nil {
		r.ctx.close()
	}
}
