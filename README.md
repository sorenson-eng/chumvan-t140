# T140 - RFC4103
___

## 1. Intro
This package includes the T140 packet implementation from the [RFC4103](https://datatracker.ietf.org/doc/html/rfc4103).
It is based on some code from [pion/rtp](https://github.com/pion/rtp), but this is not a part of pion/rtp.

## 2. T140 packet
In this implementation, **T140 packet** refers to an RTP packet whose payload is in T140 format.
A T140 packet, similar to an RTP packet, includes 2 parts:
- **T140 header** = a "constrained" RTP header
- **T140 payload** = an RTP payload contains a T140block
Examples for a packet:

Non-padding packet
valid raw packet without padding
```
0                   1                   2                   3
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|V=2|P|X| CC=0  |M|   T140 PT   |       sequence number         |
|10 |0|0| 0000  |1| 1100100(100)|  			27023	  	        | -> 0x80, 0xe4, 0x69, 0x8f
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                      timestamp (1000Hz)                       |
|				   		3653407706   					        | -> 0xd9, 0xc2, 0x93, 0xda
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|           synchronization source (SSRC) identifier            |
|					 	476325762							    | -> 0x1c, 0x64, 0x27, 0x82
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                      T.140 encoded data                       |
+                           Hello 	            		        | -> 0x48, 0x65, 0x6c, 0x6c, 0x6f
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

Packet with padding
valid raw packet with padding
```
0                   1                   2                   3
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|V=2|P|X| CC=0  |M|   T140 PT   |       sequence number         |
|10 |1|0| 0000  |1| 1100101(101)|  			27023	  	        | -> 0xa0, 0xe5, 0x69, 0x8f
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                      timestamp (1000Hz)                       |
|				   		3653407706   					        | -> 0xd9, 0xc2, 0x93, 0xda
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|           synchronization source (SSRC) identifier            |
|					 	476325762							    | -> 0x1c, 0x64, 0x27, 0x82
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                      T.140 encoded data                       |
+                           Hello 	            		        | -> 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x00, 0x03
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

```



### 2.1 T140 header
Here is the template for header:
### 2.2 T140 payload

