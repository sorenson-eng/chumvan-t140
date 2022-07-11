package t140packet

import (
	"testing"
)

func TestCountREDHeader(t *testing.T) {
	// the 1st bit is non-zero
	payload := []byte{0x00}
	if _, err := CountREDHeaders(payload); err == nil {
		t.Fatal("TestCountREDHeader: CountREDHeader did not return invalid RED header")
	}

	// Right order of the 1st bit in R-Header
	payload = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
	}
	count, err := CountREDHeaders(payload)
	if err != nil {
		t.Error(err)
	}
	if count != 3 {
		t.Errorf("TestCountREDHeader: wrong count value returned: got %#v, but want %#v", count, 3)
	}

	// counting did not stop after the R-header list
	payload = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
		0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF,
	}

	count, err = CountREDHeaders(payload)
	if err != nil {
		t.Error(err)
	}
	if count != 3 {
		if count == 5 {
			t.Errorf("TestCountREDHeader: did not stop when meet the end of R-header list")
		}
		t.Errorf("TestCountREDHeader: wrong count value returned: got %#v, but want %#v", count, 3)
	}
}

// func TestUnmarshalRHeader(t *testing.T) {
// 	t140 := &T140Packet{}
// 	// 11100101111111110000000000100000 -> 0xE5, 0xFF, 0x00, 0x20
// 	rawPacket := []byte{

// 	}
// }
