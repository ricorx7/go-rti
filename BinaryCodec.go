package rti

import (
	"bytes"
	"encoding/binary"
	"log"
	"sync"
)

var buffer bytes.Buffer

// Create some bytes to hold the header
var header [32]byte

var mutex = &sync.Mutex{}

// Ensemble ID
// 16 Bytes of 0x80 or 128
//var ensID = [...]byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80}

const hdlen = 32
const idlen = 16
const checksumSize = 4

var ensembleSize uint32
var ensembleNum uint32
var headerFound bool
var headerIndex int

// BinaryCodec decodes and passes
// all the decoded data as ensembles.
type BinaryCodec struct {
	Write     chan []byte   // Write binary data to be decoded
	Read      chan Ensemble // Read out ensembles decoded
	IsClosing bool          // Flag to stop decoding data
}

// codec creates all the channels to decod
// the data and pass it to all those that want
// the decoded ensemble data.
// var codec = binaryCodec{
// 	write:     make(chan []byte),
// 	read:      make(chan Ensemble),
// 	isClosing: false,
// }

// Run will take any incoming data and decode it.
func (codec *BinaryCodec) Run() {
	log.Print("Inside binaryCodec run")

	// buffer.Grow(4096)
	//
	// f, err := os.OpenFile("test.log", os.O_WRONLY|os.O_CREATE, 0777)
	// if err != nil {
	// 	//t.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()
	//
	// log.SetOutput(f)

	for {
		select {
		case d := <-codec.Write:
			//log.Printf("Decoding binary data: %d", len(d))

			if codec.IsClosing {
				log.Print("Closing the binary codec")
				return
			}

			mutex.Lock()

			// Add the data to the buffer
			buffer.Write(d)

			//log.Printf("AddBuffer size: %d", buffer.Len())

			// Decode the data
			decodeIncomingData()

			mutex.Unlock()
			//log.Print("End binarycodec run write")
		}
	}
}

func decodeIncomingData() {
	//log.Print("Add data to buffer")

	// Append the data to the buffer
	//buffer = append(buffer, data...)
	//buffer.Add(data)

	// If the header has not been found, find the header
	if !headerFound {

		// // Verify enough bytes are there to read the header
		//if buffer.Len() >= hdlen-headerIndex {
		// 	//temp := buffer.Bytes()
		// 	//log.Printf("buffer not big enough for header %d %d", buffer.Len(), temp[0])
		// 	//log.Printf("Header Index: %d", headerIndex)
		//log.Printf("Buffer size: %d", buffer.Len())
		//log.Printf("Header %v", header)
		// 	return
		// }

		// Look for the beginning of and ensemble
		findHeader()
		//}
	}

	// If the header was not found, return
	if !headerFound {
		return
	}

	// Ensemble size is the Header length + payload size + checksum
	var ensSize = ensembleSize + checksumSize

	// Verify enough bytes are there to read the ensemble
	if buffer.Len() < int(ensSize) {
		//log.Printf("Buffer not big enough for ensemble %d", buffer.Len())
		return
	}

	// Get the ensemble from the buffer
	payload := make([]byte, ensSize)
	size, err := buffer.Read(payload)

	// Check for error
	if err != nil {
		//log.Printf("error reading the payload: %d", size)
		return
	}

	if uint32(size) != ensSize {
		//log.Print("Ensemble not complete")
		return
	}

	// Combine the header and payload
	var ens = append(header[:], payload[:]...)

	//log.Printf("ens: %v", ens)
	//temp := buffer.Bytes()
	//log.Printf("buffer: %v", temp)

	// Get the checksum
	checksum := uint32(ens[len(ens)-4])
	checksum += uint32(ens[len(ens)-3]) << 8
	checksum += uint32(ens[len(ens)-2]) << 16
	checksum += uint32(ens[len(ens)-1]) << 24

	cal1Checksum := calculateEnsembleChecksum(ens)

	//log.Printf("Checksum: %d", checksum)
	//log.Printf("Checksum: %d  calculated1 checksum %d", checksum, cal1Checksum)

	// clear the header
	clearHeader()

	//log.Printf("Ensemble number: %d  Ensemble size: %d", ensembleNum, len(ens))
	//log.Printf("Last two bytes of ensemble: %x %x", ens[len(ens)-1], ens[len(ens)-2])

	//log.Printf("Buffer Size: %d", len(buffer))
	// Remove the ensemble from the buffer
	//buffer = buffer[ensSize:len(buffer)]
	//log.Printf("New Buffer size: %d", len(buffer))

	if uint32(cal1Checksum) == checksum {
		log.Printf("Ensemble number: %d  Ensemble size: %d", ensembleNum, len(ens))
	}

}

func findHeader() {
	//log.Printf("Search for ID %d Index: %d", buffer.Len(), headerIndex)

	// Remove the first byte until it is an ID
	for {
		// // Verify enough bytes are there to read the header
		// if buffer.Len() < hdlen {
		// 	log.Printf("Buffer not big enough to find ID %d", buffer.Len())
		// 	//log.Printf("Header Index: %d", headerIndex)
		// 	//log.Printf("Header %v", header)
		// 	return
		// }

		// Read the next byte in the buffer
		id, err := buffer.ReadByte()

		// Check for errors or if the buffer is empty
		if err != nil {
			//log.Print("Error finding ID. " + err.Error())
			return
		}

		// Check if the ID is correct
		if id != 0x80 {
			//log.Printf("Not ID: %x", id)
			//log.Printf("Header Index: %d", headerIndex)
			//if headerIndex == 1 {
			//temp := buffer.Bytes()
			//log.Printf("%v", temp[:5])
			//}
			// Start over since not an ID value
			headerIndex = 0
		} else {
			//log.Printf("Header Index: %d", headerIndex)
			header[headerIndex] = id // Set the first ID
			headerIndex++
			//log.Printf("ID: %x", id)
			//log.Printf("Header %v", header)
			//log.Printf("Header Index: %d", headerIndex)
			//log.Printf("Buffer cap: %d", buffer.Cap())

			// Found a complete ID
			if headerIndex == idlen {
				break // Break out for loop
			}
		}
	}

	//log.Print("ID found")

	//log.Print("Verify the header")

	var invEnsNum uint32
	var ensNum uint32
	var payloadSize uint32
	var invPayloadSize uint32

	// Get the next 16 bytes to get the ensemble number and payload size
	for ; headerIndex < 32; headerIndex++ {
		// Read in byte
		id, err := buffer.ReadByte()

		// Check for error
		if err != nil {
			//log.Print("Error reading ens Num and Payload size. " + err.Error())
			return
		}

		// Set the values
		header[headerIndex] = id
	}

	//log.Print("Getting Payload and Ensemble number")
	//log.Printf("Header %v", header)
	//log.Printf("header index: %d", headerIndex)

	// Look for the ensemble number and payload size
	// Get the ensemble number
	ensNum = binary.LittleEndian.Uint32(header[16:20])
	//log.Printf("Ensemble Number: %d", ensNum)

	// 1's Compliment Ensemble Number
	invEnsNum = binary.LittleEndian.Uint32(header[20:24])
	//log.Printf("Ensemble Number 1's Compliment: %d", invEnsNum)
	//log.Printf("Ensemble Number 1's Compliment: %d", ^invEnsNum)

	// Get the payload size
	payloadSize = binary.LittleEndian.Uint32(header[24:28])
	//log.Printf("Payload size: %d", payloadSize)

	// 1's Compliment Payload size
	invPayloadSize = binary.LittleEndian.Uint32(header[28:32])
	//log.Printf("Payload size 1's Compliment: %d", invPayloadSize)
	//log.Printf("Payload size Un Comliment: %d", ^invPayloadSize)

	// Verify Ensemble Number and Payload size
	if ensNum == ^invEnsNum && payloadSize == ^invPayloadSize {
		//log.Print("Payload and ensemble number matched")

		ensembleNum = ensNum
		ensembleSize = payloadSize
		headerFound = true
	}
}

// clearHeader will clear the header.
func clearHeader() {
	for i := 0; i < 32; i++ {
		header[i] = 0
	}
	headerIndex = 0
	headerFound = false

	//time.Sleep(500 * time.Millisecond)

}

/// <summary>
/// Calculate the checksum for the given ensemble.
/// This will use CRC-16 to calculate the checksum.
/// Give all bytes in the Ensemble including the checksum.
/// </summary>
/// <param name="ensemble">Byte array of the data.</param>
/// <returns>Checksum value for the given byte[].</returns>
func calculateEnsembleChecksum(ensemble []byte) uint16 {
	var crc uint16

	// Do not include the checksum to calculate the checksum
	for i := hdlen; i < len(ensemble)-checksumSize; i++ {
		crc = uint16(byte(crc>>8)) | crc<<8
		crc ^= uint16(ensemble[i])
		crc ^= uint16(byte((crc & 0xff) >> 4))
		crc ^= uint16((crc << 8) << 4)
		crc ^= uint16(((crc & 0xff) << 4) << 1)
	}

	return crc
}
