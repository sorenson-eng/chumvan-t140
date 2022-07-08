// Package t140packet implements a T140 (text/t140) packet from RFC-4103
package t140packet

import (
	"github.com/pion/rtp"
)

const payloadMaxSize uint16 = 512

// T140Packet represents a T140 packet as a form of RTP packet.
// Header is an RTP header without extensions or CSRCs.
// Payload contains redundant data if RED-flag is true.
type T140Packet struct {
	Header  rtp.Header
	Payload []byte
	IsRED   bool
}

// Unmarshal parses the passed in byte slice
// and store the result in the T140Packet this method is called upon.
// Returns any occurred error
func (t *T140Packet) Unmarshal(buf []byte, codeRED int8) (err error) {
	rtpPacket := &rtp.Packet{}
	err = rtpPacket.Unmarshal(buf)
	if err != nil {
		return
	}
	if rtpPacket.Header.Extension ||
		rtpPacket.Header.Extensions != nil ||
		rtpPacket.Header.ExtensionProfile != 0 {
		return errNoExtensionAllowed
	}
	if rtpPacket.Header.CSRC != nil {
		return errNoCSRCAllowed
	}
	if len(rtpPacket.Payload) > int(payloadMaxSize) {
		return errTooLargePayload
	}

	t.Header = rtpPacket.Header
	t.Payload = rtpPacket.Payload
	if t.Header.PayloadType == uint8(codeRED) {
		t.IsRED = true
	}
	return
}
