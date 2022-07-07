// Package t140packet implements a T140 (text/t140) packet from RFC-4103
package t140packet

// T140Packet represents a T140  in an RTP packet
type T140Packet struct {
	T140block []byte
}

// Unmarshal parses the passed in byte slice and store the result in the T140Packet this method is called upon
// Returns any error if there is
func (t *T140Packet) Unmarshal([]byte) (err error) {

	return
}
