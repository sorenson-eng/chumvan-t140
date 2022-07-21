package t140writer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewT140Writer(t *testing.T) {
	t140writer := NewT140Writer(&bytes.Buffer{})
	assert.NotNil(t, t140writer)
}
