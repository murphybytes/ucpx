package client

import (
	"github.com/murphybytes/ucp/common"
)

type writer struct {
	ctx *context
}

func newWriter(flags *common.Flags) (w common.Writer, e error) {
	var ctx *context
	ctx, e = getWriteContext(flags)
	if e != nil {
		return
	}
	w = &writer{
		ctx: ctx,
	}
	return
}

func (w *writer) Write(buffer []byte) (e error) {
	return
}

func (w *writer) Close() {
	w.ctx.close()

}
