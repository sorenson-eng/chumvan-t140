# T140 - RFC4103
*This document was made to clarify the process of adopting a payload format for RTP.*
*A better approach would be forking then merging.*
*This is used by the author to keep track of what to be done*
___

## 1. Intro
This package includes the T140 packet implementation from the [RFC4103](https://datatracker.ietf.org/doc/html/rfc4103).
It is based on some code from [pion/rtp](https://github.com/pion/rtp), but this is not a part of pion/rtp.

## 2. T140 packet
In this implementation, **T140 packet** refers to an RTP packet whose payload is in `text/t140` format.
A T140 packet, similar to an RTP packet, includes 2 parts:
- **T140 header** = a "constrained" RTP header
- **T140 payload** = an RTP payload contains a T140block

<mark>NOTE</mark> In this implementation, even `text/t140` is being considered as a type of "codec", this format does not require any further header and depends solely on the RTP header (specified in [RFC3550](https://datatracker.ietf.org/doc/html/rfc3550) with the constraints in section 2.1). The structure of this implementation may differ from the implementation in [pion/rtp/codecs](https://github.com/pion/rtp/tree/master/codecs) for this reason. For example, the `T140Packet` will be a special type embedded the `Packet` and the payload of `T140Packet` will be processed according to the `PT` field in the header (with redundancy-`RED` or an ordinary packet-`T140`).

Examples for a packet:

Non-padding packet with valid byte-slice form on the right
```
0                   1                   2                   3
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|V=2|P|X| CC=0  |M|   T140 PT   |       sequence number         |
|10 |0|0| 0000  |1| 1100100(100)|            27023              | -> 0x80, 0xe4, 0x69, 0x8f
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                      timestamp (1000Hz)                       |
|                        3653407706                             | -> 0xd9, 0xc2, 0x93, 0xda
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|           synchronization source (SSRC) identifier            |
|                          476325762                            | -> 0x1c, 0x64, 0x27, 0x82
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                      T.140 encoded data                       |
|                            Hello                              | -> 0x48, 0x65, 0x6c, 0x6c, 0x6f
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

Packet with padding with valid byte-slice form on the right
```
0                   1                   2                   3
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|V=2|P|X| CC=0  |M|   T140 PT   |       sequence number         |
|10 |1|0| 0000  |1| 1100101(100)|           27023               | -> 0xa0, 0xe5, 0x69, 0x8f
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                      timestamp (1000Hz)                       |
|                        3653407706                             | -> 0xd9, 0xc2, 0x93, 0xda
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|           synchronization source (SSRC) identifier            |
|                          476325762                            | -> 0x1c, 0x64, 0x27, 0x82
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                      T.140 encoded data                       |
|                            Hello (+ padding bytes)            | -> 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x00, 0x03
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

```

A text/t140 packet with one redundant T140block
```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|V=2|P|X| CC=0  |M|  "RED" PT   |   sequence number of primary  |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|               timestamp of primary encoding "P"               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|           synchronization source (SSRC) identifier            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|1|   T140 PT   |  timestamp offset of "R"  | "R" block length  |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|0|   T140 PT   | "R" T.140 encoded redundant data              |
+-+-+-+-+-+-+-+-+                               +---------------+
+                                               |               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+     +-+-+-+-+-+
|                "P" T.140 encoded primary data       |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

### 2.1 T140 header
From the template of RTP header, the T140-specific setting can be listed as:
- `X`: is always set to `0` which means no extension is allowed
- `CC`: is always set to `000`
- `CSRC`: due to `CC` = `000` the `CSRC` slice in the normal RTP has length 0
- `PT`: is dynamic (range from 96 to 127) and agreed via other means (e.g SIP). In this implementation, `T140 PT` is chosen to be `1100100`(100) and `RED PT`(for redundancy) is chosen to be `1100101`(101)

### 2.2 T140 payload
**The T140block**

    T.140 text is UTF-8 coded, as specified in T.140, with no extra
    framing.  The T140block contains one or more T.140 code elements as
    specified in [1].  Most T.140 code elements are single ISO 10646 [5]
    characters, but some are multiple character sequences.  Each
    character is UTF-8 encoded [6] into one or more octets.  Each block
    MUST contain an integral number of UTF-8 encoded characters
    regardless of the number of octets per character.  Any composite
    character sequence (CCS) SHOULD be placed within one block.

There is no further instruction for the payload. The payload length is variable between packets and is not determined by any means from the header.


## 3. Implementation
All components in this package is referenced from the [pion/rtp](https://github.com/pion/rtp) and [webrtc/pkg/media/](https://github.com/pion/webrtc/tree/master/pkg/media):
- [ ]`t140Packet.go`
- [ ] `t140Reader.go`
- [ ] `t140Writer.go`