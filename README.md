# T140 - RFC4103
___

## Intro
This package includes the T140 packet implementation from the [RFC4103](https://datatracker.ietf.org/doc/html/rfc4103).
It is based on some code from [pion/rtp](https://github.com/pion/rtp), but this is not a part of pion/rtp.

## T140 packet
In this implementation, **T140 packet** refers to an RTP packet whose payload is in T140 format.
A T140 packet, similar to an RTP packet, includes 2 parts:
- **T140 header** = a "constrained" RTP header
- **T140 payload** = an RTP payload contains a T140block

