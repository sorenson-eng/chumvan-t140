// Package t140packet implements a T140 (text/t140) packet from RFC-4103
package t140packet

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/pion/rtp"
)

const (
	payloadMaxSize = 128
	rHeaderSize    = 4 // bytes
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

	Payload     []byte
	PaddingSize byte
	// RED
	RHeaders []RBlockHeader

	PBlock  []byte
	RBlocks []RBlock
}

// RBlockHeader represents a header according to an "R" block.
// Includes: a payload type, a timestamp offset and a length for the "R" block
type RBlockHeader struct {
	PayloadType     uint8
	TimestampOffset uint16
	BlockLength     uint16
}

// RBlock represents an "R" block of UTF-8 data
// Includes the T140 Payload Type and the actual redundant data - "R" block
type RBlock struct {
	PayloadType uint8
	Data        []byte
}

// Unmarshal parses the passed in byte slice
// and stores the result in the T140Packet this method is called upon.
// Returns any occurred error
func (t *T140Packet) Unmarshal(buf []byte, redPT uint8) (pBlock []byte, rBlock []RBlock, err error) {
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
	if len(rtpPacket.Header.CSRC) != 0 {
		err = errNoCSRCAllowed
		return
	}
	if len(rtpPacket.Payload) > int(payloadMaxSize) {
		err = errTooLargePayload
		return
	}

	// If redundancy is applied
	t.Header = rtpPacket.Header
	if t.Header.PayloadType == redPT {
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
	// RTP payload -> T140 total payload
	t.Payload = make([]byte, len(payload))
	copy(t.Payload, payload)

	// Simple return if only P-block is in a payload
	if !t.IsRED {
		t.PBlock = payload
		return
	}

	// Payload of T140 packet can be empty
	if len(payload) == 0 {
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
	rblocks := make([]RBlock, 0)
	for _, r := range t.RHeaders {
		if r.BlockLength != 0 {
			rb := RBlock{
				PayloadType: r.PayloadType,
				Data:        payload[rLen+1 : rLen+1+int(r.BlockLength)],
			}
			rblocks = append(rblocks, rb)
			rLen += int(1 + r.BlockLength)
		}
	}
	t.RBlocks = make([]RBlock, len(t.RHeaders))
	copy(t.RBlocks, rblocks)

	t.PBlock = make([]byte, len(payload[rLen:]))
	copy(t.PBlock, payload[rLen:])
	return
}

// String returns the string representation of the T140-payload RTP packet
func (t T140Packet) String() string {
	h := t.Header
	s := "\tRTP T140 PACKET:\n"
	s += fmt.Sprintf("\tVersion: %v\n", h.Version)
	s += fmt.Sprintf("\tMarker: %v\n", h.Marker)
	s += fmt.Sprintf("\tPayload Type: %d\n", h.PayloadType)
	s += fmt.Sprintf("\tSequence Number: %d\n", h.SequenceNumber)
	s += fmt.Sprintf("\tTimestamp: %d\n", h.Timestamp)
	s += fmt.Sprintf("\tSSRC: %d (%x)\n", h.SSRC, h.SSRC)
	s += fmt.Sprintf("\tCSRC: %v\n", h.CSRC)
	s += fmt.Sprintf("\tIs RED: %t\n", t.IsRED)
	s += fmt.Sprintf("\tP-block length: %d bytes\n", len(t.PBlock))
	s += fmt.Sprintf("\tR-blocks quantity: %d\n", len(t.RBlocks))
	s += fmt.Sprintf("\tR-blocks: %v\n", t.RBlocks)
	s += fmt.Sprintf("\tPayload: %v\n", t.Payload)
	return s
}

//T140Payloader payloads T140 packets
type T140Payloader struct{}

// Payload fragments a packet across one or more byte array.
// In T140 packet, each being sent data is constrained to 1 packet.
// The operation still return a 2-dimensional byte slice to conform with the interface.
func (p *T140Payloader) Payload(mtu uint16, payload []byte) (payloads [][]byte, err error) {
	if payload == nil {
		return
	}

	if len(payload) > payloadMaxSize {
		return payloads, errTooLargePayload
	}

	out := make([]byte, len(payload))
	copy(out, payload)
	payloads = [][]byte{out}

	return
}

// Marshal serializes the calling packet to a newly created slice of bytes.
// Returns the serialized bytes and any occurred error.
func (t T140Packet) Marshal() (buf []byte, err error) {
	buf = make([]byte, t.MarshalSize())

	n, err := t.MarshalTo(buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}

// MarshalTo takes in a byte slice, marshals the calling packet into that slice.
// Returns number of bytes being marshaled and any occurred error
func (t T140Packet) MarshalTo(buf []byte) (n int, err error) {
	n, err = t.Header.MarshalTo(buf)
	if err != nil {
		return 0, err
	}
	// Make sure the buffer is large enough to hold the packet.
	if n+len(t.Payload)+int(t.PaddingSize) > len(buf) {
		return 0, io.ErrShortBuffer
	}

	m := copy(buf[n:], t.Payload)
	if t.Header.Padding {
		buf[n+m+int(t.PaddingSize-1)] = t.PaddingSize
	}

	return n + m + int(t.PaddingSize), nil
}

// MarshalSize returns the size of the packet once marshaled.
func (t T140Packet) MarshalSize() int {
	return t.Header.MarshalSize() + len(t.Payload) + int(t.PaddingSize)
}

// ToRTP returns an RTP packet based on the calling T140 packet
func (t T140Packet) ToRTP() (r *rtp.Packet) {
	r = &rtp.Packet{}
	r.Header = t.Header.Clone()
	if t.Payload != nil {
		r.Payload = make([]byte, len(t.Payload))
		copy(r.Payload, t.Payload)
	}
	r.PaddingSize = t.PaddingSize
	return r
}
