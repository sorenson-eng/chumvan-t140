package t140writer

import (
	"io"

	"github.com/pion/rtpio/pkg/rtpio"
)

type T140Writer interface {
	rtpio.RTPWriter
}

func NewT140Writer(w io.Writer) T140Writer {
	return rtpio.NewRTPWriter(w)
}
