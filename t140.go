package t140

import (
	"io"

	"github.com/chumvan/t140/t140writer"
)

type T140WriterCloser interface {
	t140writer.T140Writer
	io.Closer
}
