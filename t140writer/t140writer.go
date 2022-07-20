package t140writer

import (
	"errors"

	"github.com/pion/rtpio/pkg/rtpio"
)

type T140Writer struct {
	rtpio.RTPWriter
}

type T140WriteCloser interface {
	T140Writer
	Close() error
}

var errClose = errors.New("Close error")

func (t *T140Writer) Close() error {
	return errClose
}
