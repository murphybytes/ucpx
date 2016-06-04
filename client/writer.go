package client

import (
	"io"

	"github.com/murphybytes/ucp/common"
)

func newWriter(flags *common.Flags) (w io.WriteCloser, e error) {
	return getWriteContext(flags)
}
