package t140packet

import (
	"testing"
)

func TestCountREDHeader(t *testing.T) {
	payload := []byte{0x00, 0x00}
	if _, err := CountREDHeaders(payload); err == nil {
		t.Fatal("CountREDHeader did not return invalid RED header")
	}
}
