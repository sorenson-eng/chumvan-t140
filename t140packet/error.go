package t140packet

import "errors"

var (
	errNoExtensionAllowed        = errors.New("no extension is allowed in T140 packet")
	errNoCSRCAllowed             = errors.New("CC must be 0 and CSRC must be nil in T140 packet")
	errTooLargePayload           = errors.New("payload is too large for T140 packet, keep it short")
	errInvalidREDHeader          = errors.New("redundant block header must start with bit 1")
	errMismatchREDHeaderAndBlock = errors.New("mismatch in number of RED header and data block")
	errMismatchRHeaderSize       = errors.New("redundant block header size must be 4 bytes")
)
