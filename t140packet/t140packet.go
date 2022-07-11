// Package t140packet implements a T140 (text/t140) packet from RFC-4103
package t140packet

import (
	"encoding/binary"

	"github.com/pion/rtp"
)

const (
	payloadMaxSize = 512
	rHeaderSize    = 4 // byte
	rHeaderMask    = 0x80
	ptMask         = 0x7F

	timeOffsetShift = 2
	timeOffsetSize  = 2 // 14-bit

	rBlockLengthMask = 0x03FF
	rBlockLengthSize = 2 // 10-bit
)

// T140Packet represents a T140 packet as a form of RTP packet.
// Header is an RTP header without extensions or CSRCs.
// Payload contains redundant data if RED-flag is true.
type T140Packet struct {
	Header rtp.Header
	IsRED  bool

	// RED
	RHeaders []RBlockHeader

	PBlock  []byte
	RBlocks []RBlock
}

type RBlockHeader struct {
	PayloadType     uint8
	TimestampOffset uint16
	BlockLength     uint16
}

type RBlock struct {
	PayloadType uint8
	Data        []byte
}

// Unmarshal parses the passed in byte slice
// and stores the result in the T140Packet this method is called upon.
// Returns any occurred error
func (t *T140Packet) Unmarshal(buf []byte, codeRED uint8) (pBlock []byte, rBlock []RBlock, err error) {
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

	err = t.UnmarshalPayload(rtpPacket.Payload)
	if err != nil {
		return
	}

	return
}

// UnmarshalBlock parses the passed in byte slice
// and stores the block(s) in the T140Packet this method is called upon.
// Returns any occurred error.
func (t *T140Packet) UnmarshalPayload(payload []byte) (err error) {
	// Payload of T140 packet can be empty
	if len(payload) == 0 {
		return
	}
	// Simple return if only P-block is in a payload
	if !t.IsRED {
		t.PBlock = payload
		return
	}

	// Handle multi-blocks (with redundancy) in a payload
	err = t.UnmarshalRHeaders(payload)
	if err != nil {
		return
	}

	err = t.unmarshalBlocks(payload)
	if err != nil {
		return
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
	// TODO check out of order RHeaders
	rowCount := len(payload) / rHeaderSize
	for i := 0; i <= rowCount; i++ {
		if payload[i*rHeaderSize]&rHeaderMask == 0x80 {
			count++
		} else {
			return
		}
	}

	return
}

// UnmarshalRHeader parses the passed in byte slice
// and stores parsed RHeaders into the T140 packet this method is called upon.
// Returns any occurred error
func (t *T140Packet) UnmarshalRHeaders(payload []byte) (err error) {
	rCount, err := CountREDHeaders(payload)
	if err != nil {
		return err
	}

	for i := 0; i < rCount; i++ {
		buf := make([]byte, rHeaderSize)
		copy(buf, payload[i*rHeaderSize:(i+1)*rHeaderSize])
		if buf[0]&rHeaderMask == 0 {
			return errInvalidREDHeader
		}
		rHeader := &RBlockHeader{}
		rHeader.PayloadType = buf[0] & ptMask
		rHeader.TimestampOffset = binary.BigEndian.Uint16(buf[1:1+timeOffsetSize]) >> timeOffsetShift
		rHeader.BlockLength = binary.BigEndian.Uint16(buf[2:]) & uint16(rBlockLengthMask)

		t.RHeaders = append(t.RHeaders, *rHeader)
	}
	return
}

// unmarshalBlocks parses the passed in byte slice
// and stores the "R" block and "P" block in the T140 packet
// this method is called upon.
// Returns any occurred error
func (t *T140Packet) unmarshalBlocks(payload []byte) (err error) {
	var rLen int = len(t.RHeaders) * rHeaderSize
	for _, r := range t.RHeaders {
		if r.BlockLength != 0 {
			rb := RBlock{
				PayloadType: r.PayloadType,
				Data:        payload[rLen+1 : rLen+1+int(r.BlockLength)],
			}
			t.RBlocks = append(t.RBlocks, rb)
			rLen += int(1 + r.BlockLength)
		}
	}

	t.PBlock = payload[rLen:]
	return
}
