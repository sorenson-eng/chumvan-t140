package t140packet

import (
	"reflect"
	"testing"

	"github.com/pion/rtp"
)

func TestCountREDHeader(t *testing.T) {
	// the 1st bit is non-zero
	payload := []byte{0x00}
	if _, err := CountREDHeaders(payload); err == nil {
		t.Fatal("TestCountREDHeader did not return invalid RED header")
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
		t.Errorf("TestCountREDHeader wrong count value returned: got %#v, but want %#v", count, 3)
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
			t.Errorf("TestCountREDHeader did not stop when meet the end of R-header list")
		}
		t.Errorf("TestCountREDHeader wrong count value returned: got %#v, but expect %#v", count, 3)
	}
}

func TestUnmarshalRHeaders(t *testing.T) {
	t140 := &T140Packet{}
	// 11100101 11111111 00000000 00001010
	payload := []byte{
		0xE4, 0xFF, 0x00, 0x0A, // RHeader
		0x64,                                                       // 0-flag and T140 PT
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, // "R" data
		0x01, 0x02, 0x03, 0x04, // "P" data
	}
	if err := t140.UnmarshalRHeaders(payload); err != nil {
		t.Error(err)
	}

	expectRHeader := RBlockHeader{
		PayloadType:     100,
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
		0xE4, 0xFF, 0x00, 0x0A, // RHeader
		0xE4, 0xFF, 0x00, 0x0A, // RHeader
		0x64,                                                       // 0-flag and T140 PT
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, // "R" data
		0x64,                                                       // 0-flag and T140 PT
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, // "R" data
		0x01, 0x02, 0x03, 0x04, // "P" data
	}
	if err := t140.UnmarshalRHeaders(payload); err != nil {
		t.Error(err)
	}

	if len(t140.RHeaders) != 2 {
		t.Errorf("TestUnmarshalRHeaders wrong number of RHeaders return: got %d, but expect: %d", len(t140.RHeaders), 2)
	}

	expectRHeaders := []RBlockHeader{
		{
			PayloadType:     100,
			TimestampOffset: 16320,
			BlockLength:     10,
		},
		{
			PayloadType:     100,
			TimestampOffset: 16320,
			BlockLength:     10,
		},
	}
	if !reflect.DeepEqual(t140.RHeaders, expectRHeaders) {
		t.Errorf("TestUnmarshalHeaders mismatch unmarshal RHeader: got %#v, but expect %#v", firstRHeader, expectRHeader)
	}
}

func TestUnmarshalBlocks(t *testing.T) {
	t140 := &T140Packet{}
	// empty payload
	if err := t140.unmarshalBlocks([]byte{}); err != nil {
		t.Errorf("Test UnmarshalBlocks empty payload should be allowed")
	}

	// P-block only
	payload := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
	t140.IsRED = false
	if err := t140.unmarshalBlocks(payload); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(t140.PBlock, payload) {
		t.Errorf("TestUnmarshalBlocks incorrect unmarshal payload: got %#v, but expect: %#v", t140.PBlock, payload)
	}

	// with R-block
	payload = []byte{
		0xE4, 0xFF, 0x00, 0x0A, // RHeader
		0xE4, 0xFF, 0x00, 0x0A, // RHeader
		0x64,                                                       // 0-flag and T140 PT
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, // "R" data
		0x64,                                                       // 0-flag and T140 PT
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, // "R" data
		0x01, 0x02, 0x03, 0x04, // "P" data
	}
	RBlocks := []RBlock{
		{
			PayloadType: 100,
			Data:        payload[9:19],
		},
		{
			PayloadType: 100,
			Data:        payload[9:19],
		},
	}
	t140 = &T140Packet{}
	t140.IsRED = true
	if err := t140.UnmarshalRHeaders(payload); err != nil {
		t.Error(err)
	}
	if err := t140.unmarshalBlocks(payload); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(t140.PBlock, payload[30:34]) {
		t.Errorf("TestUnmarshalBlocks incorrect unmarshal P block: got %#v, but expect: %#v", t140.PBlock, payload[30:34])
	}
	if !reflect.DeepEqual(t140.RBlocks, RBlocks) {
		t.Errorf("TestUnmarshalBlocks incorrect unmarshal R blocks: got %#v, but expect: %#v", t140.RBlocks, RBlocks)
	}
}

func TestUnmarshal(t *testing.T) {
	var redPT uint8 = 101
	var t140PT uint8 = 100
	// only P-block
	t140 := &T140Packet{}
	rawPacket := []byte{
		0x80, 0xe4, 0x69, 0x8f,
		0xd9, 0xc2, 0x93, 0xda,
		0x1c, 0x64, 0x27, 0x82,
		0x48, 0x65, 0x6c, 0x6c, 0x6f,
	}

	marshalPacket := &T140Packet{
		Header: rtp.Header{
			Version:        2,
			Padding:        false,
			Extension:      false,
			Marker:         true,
			PayloadType:    100,
			SequenceNumber: 27023,
			Timestamp:      3653407706,
			SSRC:           476325762,
			CSRC:           []uint32{},
		},
		PBlock: rawPacket[12:],
	}
	if _, _, err := t140.Unmarshal(rawPacket, redPT); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(t140, marshalPacket) {
		t.Errorf("TestUnmarshal incorrect unmarshal packet: got %#v, but expect %#v", t140, marshalPacket)
	}

	// with RBlocks
	t140 = &T140Packet{}
	rawPacket = []byte{
		0x80, 0xe5, 0x69, 0x8f, //	---------
		0xd9, 0xc2, 0x93, 0xda, // 	RTP Header
		0x1c, 0x64, 0x27, 0x82, //	---------
		0xe4, 0xff, 0x00, 0x0a, // "R" Block Header
		0xe4, 0xff, 0x00, 0x0a, // "R" Block Header
		0x64,                                                       // 0-flag and T140 PT
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, // "R" data
		0x64,                                                       // 0-flag and T140 PT
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, // "R" data
		0x48, 0x65, 0x6c, 0x6c, 0x6f, // "P" data
	}
	marshalPacket.Header.PayloadType = redPT
	marshalPacket.IsRED = true
	marshalPacket.RBlocks = []RBlock{
		{
			PayloadType: t140PT,
			Data:        rawPacket[21:31],
		},
		{
			PayloadType: t140PT,
			Data:        rawPacket[32:42],
		},
	}
	marshalPacket.RHeaders = []RBlockHeader{
		{
			PayloadType:     100,
			TimestampOffset: 16320,
			BlockLength:     10,
		},
		{
			PayloadType:     100,
			TimestampOffset: 16320,
			BlockLength:     10,
		},
	}
	if _, _, err := t140.Unmarshal(rawPacket, redPT); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(t140, marshalPacket) {
		t.Errorf("TestUnmarshal incorrect unmarshal packet: got %#v,\n but expect %#v", t140, marshalPacket)
	}
}

func TestT140Payloader(t *testing.T) {
	var mtu uint16 = 56

	// empty-payload
	payload := []byte{}
	p := &T140Payloader{}
	payloads, err := p.Payload(mtu, payload)
	if len(payloads[0]) != 0 || err != nil {
		t.Errorf("TestT140Payloader empty payload should result in empty list of payloads: got %#v of length %d", payloads, len(payloads))
	}

	// Oversize payload
	oversize := 129
	payload = make([]byte, oversize)
	for i := 0; i < oversize; i++ {
		payload[i] = 0x01
	}
	p = &T140Payloader{}
	payloads, err = p.Payload(mtu, payload)
	if payloads != nil || err != errTooLargePayload {
		t.Errorf("TestT140Payloader oversize payload should return empty array and errTooLargePayload - error: got %v and %v", payloads, err)
	}

	// valid payload
	payload = make([]byte, 128)
	for i := range payload {
		payload[i] = 0x01
	}
	p = &T140Payloader{}
	payloads, err = p.Payload(mtu, payload)
	if err != nil || !reflect.DeepEqual(payloads[0], payload) {
		t.Errorf("TestT140Payloader max-size payload : got %#v of %T and %v\n, but want %#v of %T and %v", payloads[0], payloads[0], err, payload, payload, nil)
	}
}
