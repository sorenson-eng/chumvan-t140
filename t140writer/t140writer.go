package t140writer

import (
	"io"

	"github.com/pion/rtpio/pkg/rtpio"
)

type T140Writer interface {
	rtpio.RTPWriter
}

type T140WriteCloser interface {
	rtpio.RTPWriter
	io.Closer
}
