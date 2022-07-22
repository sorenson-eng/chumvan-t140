package t140writer

import (
	"bytes"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewT140Writer(t *testing.T) {
	t140writer := NewT140Writer(&bytes.Buffer{})
	assert.NotNil(t, t140writer)
}

func TestNewT140WriteCloser(t *testing.T) {
	receiverAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", "127.0.0.1", 6420))
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp4", nil, receiverAddr)
	if err != nil {
		panic(err)
	}
	t140writeCloser := NewT140Writer(conn)
	assert.NotNil(t, t140writeCloser)
}
