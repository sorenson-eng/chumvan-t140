package t140packet

import (
	"reflect"
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
		t.Errorf("TestCountREDHeader: wrong count value returned: got %#v, but expect %#v", count, 3)
	}
}

func TestUnmarshalRHeaders(t *testing.T) {
	t140 := &T140Packet{}
	// 11100101 11111111 00000000 00001010
	payload := []byte{
		0xE5, 0xFF, 0x00, 0x0A, // RHeader
		0x65,                                                       // 0-flag and T140 PT
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, // "R" data
		0x01, 0x02, 0x03, 0x04, // "P" data
	}
	err := t140.UnmarshalRHeaders(payload)
	if err != nil {
		t.Error(err)
	}

	expectRHeader := RBlockHeader{
		PayloadType:     101,
		TimestampOffset: 16320,
		BlockLength:     10,
	}

	firstRHeader := t140.RHeaders[0]
	if !reflect.DeepEqual(firstRHeader, expectRHeader) {
		t.Errorf("TestUnmarshalHeaders mismatch unmarshal RHeader: got %#v, but expect %#v", firstRHeader, expectRHeader)
	}

	// multiple (2) RHeaders
	t140 = &T140Packet{}
	payload = []byte{
		0xE5, 0xFF, 0x00, 0x0A, // RHeader
		0xE5, 0xFF, 0x00, 0x0A, // RHeader
		0x65,                                                       // 0-flag and T140 PT
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, // "R" data
		0x65,                                                       // 0-flag and T140 PT
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, // "R" data
		0x01, 0x02, 0x03, 0x04, // "P" data
	}
	err = t140.UnmarshalRHeaders(payload)
	if err != nil {
		t.Error(err)
	}
	if len(t140.RHeaders) != 2 {
		t.Errorf("TestUnmarshalRHeaders wrong number of RHeaders return: got %d, but expect: %d", len(t140.RHeaders), 2)
	}

	expectRHeaders := []RBlockHeader{
		{
			PayloadType:     101,
			TimestampOffset: 16320,
			BlockLength:     10,
		},
		{
			PayloadType:     101,
			TimestampOffset: 16320,
			BlockLength:     10,
		},
	}
	if !reflect.DeepEqual(t140.RHeaders, expectRHeaders) {
		t.Errorf("TestUnmarshalHeaders mismatch unmarshal RHeader: got %#v, but expect %#v", firstRHeader, expectRHeader)
	}
}
