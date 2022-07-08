// Package t140packet implements a T140 (text/t140) packet from RFC-4103
package t140packet

import (
	"encoding/binary"

	"github.com/pion/rtp"
)

const (
	payloadMaxSize = 512
	rHeaderSize    = 4
	rHeaderMask    = 0x08
	ptMask         = 0x7F

	timeOffsetOffset = 1
	timeOffsetShift  = 2
	timeOffsetMask   = 0x3FFF
	timeOffsetSize   = 2 // 14 bit
)

// T140Packet represents a T140 packet as a form of RTP packet.
// Header is an RTP header without extensions or CSRCs.
// Payload contains redundant data if RED-flag is true.
type T140Packet struct {
	Header rtp.Header
	IsRED  bool

	// RED
	RHeader []RBlockHeader

	PBlock []byte
	RBlock []byte
}

type RBlockHeader struct {
	PayloadType     uint8
	TimestampOffset uint16
	BlockLength     uint16
}

// Unmarshal parses the passed in byte slice
// and stores the result in the T140Packet this method is called upon.
// Returns any occurred error
func (t *T140Packet) Unmarshal(buf []byte, codeRED uint8) (pBlock []byte, rBlock []byte, err error) {
	rtpPacket := &rtp.Packet{}
	err = rtpPacket.Unmarshal(buf)
	if err != nil {
		return
	}
	if rtpPacket.Header.Extension ||
		rtpPacket.Header.Extensions != nil ||
		rtpPacket.Header.ExtensionProfile != 0 {
		err = errNoExtensionAllowed
		return
	}
	if rtpPacket.Header.CSRC != nil {
		err = errNoCSRCAllowed
		return
	}
	if len(rtpPacket.Payload) > int(payloadMaxSize) {
		err = errTooLargePayload
		return
	}

	// If redundancy is applied
	t.Header = rtpPacket.Header
	if t.Header.PayloadType == uint8(codeRED) {
		t.IsRED = true
	}

	err = t.UnmarshalBlock(rtpPacket.Payload)
	if err != nil {
		return
	}

	return
}

// UnmarshalBlock parses the passed in byte slice
// and stores the block(s) in the T140Packet this method is called upon.
// Returns any occurred error.
func (t *T140Packet) UnmarshalBlock(payload []byte) (err error) {
	// Payload of T140 packet can be empty
	if len(payload) == 0 {
		return
	}

	if !t.IsRED {
		t.PBlock = payload
		return
	}

	rCount, err := CountREDHeaders(payload)
	if err != nil {
		return
	}

	for i := 0; i <= rCount; i++ {
		payload[i*rHeaderSize]
	}

	return
}

// CountREDHeaders checks and counts the number of RED headers in the passed in byte slice.
// Returns number of RED headers and any occurred error.
func CountREDHeaders(payload []byte) (count int, err error) {
	if payload[0]&rHeaderMask == 0 {
		err = errInvalidREDHeader
		return
	}

	rowCount := len(payload) / rHeaderSize
	for i := 0; i <= rowCount; i++ {
		if payload[i*rHeaderSize]&rHeaderMask == 1 {
			count++
		}
	}

	return
}

func (t *T140Packet) UnmarshalRHeader(buf []byte) (err error) {
	if len(buf) != 4 {
		return errMismatchRHeaderSize
	}
	if buf[0]&rHeaderMask == 0 {
		return errInvalidREDHeader
	}
	rHeader := &RBlockHeader{}
	rHeader.PayloadType = buf[0] & ptMask
	rHeader.TimestampOffset = binary.BigEndian.Uint16(buf[timeOffsetOffset : timeOffsetOffset+tim])
}
