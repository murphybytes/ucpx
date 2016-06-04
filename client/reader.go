package client

import (
	"io"

	"github.com/murphybytes/ucp/common"
)

func newReader(flags *common.Flags) (r io.ReadCloser, e error) {

	return getReadContext(flags)

}
