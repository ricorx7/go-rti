package rti

import (
	"bytes"
	"encoding/binary"
	"log"
)

var buffer []byte

// Ensemble ID
// 16 Bytes of 0x80 or 128
var ensID = [...]byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80}

// binaryCodec decodes and passes
// all the decoded data as ensembles.
type binaryCodec struct {
	write     chan []byte   // Write binary data to be decoded
	read      chan Ensemble // Read out ensembles decoded
	isClosing bool          // Flag to stop decoding data
}

// codec creates all the channels to decod
// the data and pass it to all those that want
// the decoded ensemble data.
var codec = binaryCodec{
	write:     make(chan []byte),
	read:      make(chan Ensemble),
	isClosing: false,
}

// run will take any incoming data and decode it.
func (codec *binaryCodec) run() {
	log.Print("Inside binaryCodec run")

	for {
		select {
		case d := <-codec.write:
			log.Printf("Decoding binary data: %d", len(d))

			// Decode the data
			decodeIncomingData(d)

			if codec.isClosing {
				log.Print("Closing the binary codec")
				return
			}
		}
	}
}

func decodeIncomingData(data []byte) {
	// Append the data to the buffer
	buffer = append(buffer, data...)

	// Remove the first byte until it is an ID
	for {
		if len(buffer) > 0 {
			if buffer[0] != 0x80 {
				log.Printf("%d", buffer[0])
				// Pop from the front
				_, buffer = buffer[len(buffer)-1], buffer[:len(buffer)-1]
			} else {
				// Break out
				break
			}

		} else {
			// No more data in the buffer
			return
		}
	}

	// Look for the ID
	if len(buffer) > 32 {
		if buffer[0] == 0x80 &&
			buffer[1] == 0x80 &&
			buffer[2] == 0x80 &&
			buffer[3] == 0x80 &&
			buffer[4] == 0x80 &&
			buffer[5] == 0x80 &&
			buffer[6] == 0x80 &&
			buffer[7] == 0x80 &&
			buffer[8] == 0x80 &&
			buffer[9] == 0x80 &&
			buffer[10] == 0x80 &&
			buffer[11] == 0x80 &&
			buffer[12] == 0x80 &&
			buffer[13] == 0x80 &&
			buffer[14] == 0x80 &&
			buffer[15] == 0x80 {
			// Look for the ensemble number and payload size
			// Get the ensemble number
			en := bytes.NewBuffer(buffer[16:19])
			var ensNum uint32
			binary.Read(en, binary.LittleEndian, &ensNum)
			log.Printf("Ensemble Number: %d", ensNum)

			// 1's Compliment Ensemble Number
			ien := bytes.NewBuffer(buffer[20:23])
			var invEnsNum uint32
			binary.Read(ien, binary.LittleEndian, &invEnsNum)
			log.Printf("Ensemble Number 1's Compliment: %d", invEnsNum)

			// Get the payload size
			ps := bytes.NewBuffer(buffer[24:27])
			var payloadSize uint32
			binary.Read(ps, binary.LittleEndian, &payloadSize)
			log.Printf("Payload size: %d", payloadSize)

			// 1's Compliment Payload size
			ips := bytes.NewBuffer(buffer[28:31])
			var invPayloadSize uint32
			binary.Read(ips, binary.LittleEndian, &invPayloadSize)
			log.Printf("Payload size 1's Compliment: %d", invPayloadSize)

			// Verify Ensemble Number and Payload size

			// Use the Payload size to get the checksum value

		} else {
			// ID not found, wait for the next block of data
			return
		}
	} else {
		// ID not found, wait for the next block of data
		return
	}
}
