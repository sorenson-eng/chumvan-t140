package t140writer

import (
	"bytes"
	"testing"

	"github.com/pion/rtpio/pkg/rtpio"
	"github.com/stretchr/testify/assert"
)

func TestNewWith(t *testing.T) {
	rtpwriter := rtpio.NewRTPWriter(&bytes.Buffer{})
	w := &T140Writer{rtpwriter}
	assert.NotNil(t, w.Close())
}
