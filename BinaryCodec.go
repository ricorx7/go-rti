package rti

import (
	"bytes"
	"encoding/binary"
	"log"
	"runtime/debug"
	"strings"
	"sync"
)

/*
 * HISTORY
 * -----------------------------------------------------------------
 * Date            Initials    Version    Comments
 * -----------------------------------------------------------------
 * 09/08/2015      RC          1.0        Initial coding
 *
 *
 */

var buffer bytes.Buffer   // Buffer the incoming data
var header [32]byte       // Create some bytes to hold the header
var mutex = &sync.Mutex{} // mutex to prevent writing to the buffer while using it

// Ensemble ID
// 16 Bytes of 0x80 or 128
//var ensID = [...]byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80}

const hdlen = 32       // Header Length
const idlen = 16       // ID Length
const checksumSize = 4 // Checksum size

const dataTypeFloat = 10 // Default value for the base data type FLOAT (ValueType)
const dataTypeInt = 20   // Default value for the base data type INTEGER (ValueType)
const dataTypeByte = 50  // Default value for the base data type BYTE (ValueType)

const defaultImag = 0              // Default value for IMG in base data
const defaultNameLength = 8        // Default value for the Name Length in the base data
const numDataSetHeaderElements = 6 // This is number of elements in the header of a dataset.  Each element is a byte except the NAME.  Its size varies and is given by NameLength.
const payloadHeaderLen = 28        // Each payload contains a header with DataType, Bins or Elements, Beams or 1, Image, ID (Name) Length and ID (name)

const beamVelocityID = "E000001"       // Beam Velocity Dataset ID
const instrumentVelocityID = "E000002" // Instrument Velocity Dataset ID
const earthVelocityID = "E000003"      // Earth Velocity Dataset ID
const amplitudeID = "E000004"          // Amplitude Dataset ID
const correlationID = "E000005"        // Correlation Dataset ID
const goodBeamID = "E000006"           // Good Beam Dataset ID
const goodEarthID = "E000007"          // Good Earth Dataset ID
const ensembleDataID = "E000008"       // Ensemble Dataset ID
const ancillaryID = "E000009"          // Ancillary Dataset ID
const bottomTrackID = "E000010"        // Bottom Track Dataset ID

var ensembleSize uint32 // Current ensemble size
var ensembleNum uint32  // Current ensemble Number
var headerFound bool    // Flag if header is found
var headerIndex int     // Header index

// BinaryCodec decodes and passes
// all the decoded data as ensembles.
type BinaryCodec struct {
	Write          chan []byte   // Write binary data to be decoded
	Read           chan Ensemble // Read out ensembles decoded
	IsClosing      bool          // Flag to stop decoding data
	bufferIncoming chan []byte   // Buffer the incoming data to decode
}

// Init will initialize the codec.
func (codec *BinaryCodec) Init() {
	codec.bufferIncoming = make(chan []byte, 1024)
	go codec.Run()
}

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
			//log.Print("Start write data to codec")
			//log.Printf("Decoding binary data: %d", len(d))

			if codec.IsClosing {
				log.Print("Closing the binary codec")
				return
			}

			//log.Print("End binarycodec run write")
			select {
			case codec.bufferIncoming <- d: // Buffer the data to decode
			default:
				log.Print("Error writing data to buffer")
			}
			//log.Print("Stop write data to codec")
		case d := <-codec.bufferIncoming:
			//log.Print("Start write data to buffer")
			//mutex.Lock()

			// Add the data to the buffer
			_, err := buffer.Write(d)
			if err != nil {
				log.Print(err)
			}

			//log.Printf("AddBuffer size: %d", buffer.Len())
			//log.Printf("AddBuffer size: %d", len(d))

			// Decode the data
			decodeIncomingData(codec)

			//mutex.Unlock()
			//log.Print("Stop write data to buffer")
		default:
			//log.Print("Error in codec")
		}
	}
}

func decodeIncomingData(codec *BinaryCodec) {
	//log.Print("Add data to buffer")

	// If the header has not been found, find the header
	if !headerFound {
		// Look for the beginning of and ensemble
		findHeader()
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

	log.Printf("Checksum: %d", checksum)
	//log.Printf("Checksum: %d  calculated1 checksum %d", checksum, cal1Checksum)

	// Clear the header
	clearHeader()

	//log.Printf("Ensemble number: %d  Ensemble size: %d", ensembleNum, len(ens))
	//log.Printf("Last two bytes of ensemble: %x %x", ens[len(ens)-1], ens[len(ens)-2])

	if uint32(cal1Checksum) == checksum {
		//log.Printf("Ensemble number: %d  Ensemble size: %d", ensembleNum, len(ens))
		// Decode the ensemble
		ensemble := decodeEnsemble(ens)

		log.Print("Ensemble decoded")

		// Publish the ensemble
		codec.Read <- ensemble

		log.Printf("Ensemble number: %d Number of Beams: %d  Number of Bins: %d", ensemble.EnsembleData.EnsembleNumber, ensemble.EnsembleData.NumBeams, ensemble.EnsembleData.NumBins)
		//
		// ensJSON, _ := json.Marshal(ensemble)
		//
		// fmt.Println(string(ensJSON))

	}

	// Garbage collect
	debug.FreeOSMemory()
}

func decodeEnsemble(data []byte) Ensemble {
	// Keep track where in the packet
	// we are currently decoding
	var packetPointer = hdlen
	var enstype uint32
	var numElements uint32
	var elementMultiplier uint32
	var imag uint32
	var nameLen uint32
	var name string
	var dataSetSize int
	var ensemble Ensemble

	for i := 0; i < MaxNumDataSets; i++ {
		// Ensemble type
		ptr := packetPointer + (BytesInInt32 * 0)
		enstype = binary.LittleEndian.Uint32(data[ptr : ptr+4])
		//log.Printf("Ensemble Type: %d", enstype)

		// Number of elements
		ptr = packetPointer + (BytesInInt32 * 1)
		numElements = binary.LittleEndian.Uint32(data[ptr : ptr+4])
		//log.Printf("Num Elements: %d", numElements)

		// Element Mulitiplier
		ptr = packetPointer + (BytesInInt32 * 2)
		elementMultiplier = binary.LittleEndian.Uint32(data[ptr : ptr+4])
		//log.Printf("Element Multiplier: %d", elementMultiplier)

		// Image
		ptr = packetPointer + (BytesInInt32 * 3)
		imag = binary.LittleEndian.Uint32(data[ptr : ptr+4])
		//log.Printf("Image: %d", imag)

		// Name Length
		ptr = packetPointer + (BytesInInt32 * 4)
		nameLen = binary.LittleEndian.Uint32(data[ptr : ptr+4])
		//log.Printf("Name Length: %d", nameLen)

		// DataSet Name
		ptr = packetPointer + (BytesInFloat * 5)
		name = string(data[ptr : ptr+8])
		//log.Printf("Name: %s", name)

		// Data set size
		dataSetSize = getDataSetSize(enstype, nameLen, numElements, elementMultiplier)

		// Beam Velocity Dataset
		if strings.Contains(name, ensembleDataID) {
			// Base
			ensemble.EnsembleData.Base.Enstype = enstype
			ensemble.EnsembleData.Base.NumElements = numElements
			ensemble.EnsembleData.Base.ElementMultiplier = elementMultiplier
			ensemble.EnsembleData.Base.Imag = imag
			ensemble.EnsembleData.Base.NameLen = nameLen
			ensemble.EnsembleData.Base.Name = name

			// Ensemble Data Set
			ensemble.EnsembleData.Decode(data[packetPointer : packetPointer+dataSetSize])

			// Move to the next dataset
			packetPointer += dataSetSize
		} else if strings.Contains(name, ancillaryID) {
			// Base
			ensemble.AncillaryData.Base.Enstype = enstype
			ensemble.AncillaryData.Base.NumElements = numElements
			ensemble.AncillaryData.Base.ElementMultiplier = elementMultiplier
			ensemble.AncillaryData.Base.Imag = imag
			ensemble.AncillaryData.Base.NameLen = nameLen
			ensemble.AncillaryData.Base.Name = name

			// Ancillary Data Set
			ensemble.AncillaryData.Decode(data[packetPointer : packetPointer+dataSetSize])

			// Move to the next dataset
			packetPointer += dataSetSize
		} else if strings.Contains(name, beamVelocityID) {
			// Base
			ensemble.BeamVelocityData.Base.Enstype = enstype
			ensemble.BeamVelocityData.Base.NumElements = numElements
			ensemble.BeamVelocityData.Base.ElementMultiplier = elementMultiplier
			ensemble.BeamVelocityData.Base.Imag = imag
			ensemble.BeamVelocityData.Base.NameLen = nameLen
			ensemble.BeamVelocityData.Base.Name = name

			// Ensemble Data Set
			ensemble.BeamVelocityData.Decode(data[packetPointer : packetPointer+dataSetSize])

			// Move to the next dataset
			packetPointer += dataSetSize
		} else if strings.Contains(name, instrumentVelocityID) {
			// Base
			ensemble.InstrumentVelocityData.Base.Enstype = enstype
			ensemble.InstrumentVelocityData.Base.NumElements = numElements
			ensemble.InstrumentVelocityData.Base.ElementMultiplier = elementMultiplier
			ensemble.InstrumentVelocityData.Base.Imag = imag
			ensemble.InstrumentVelocityData.Base.NameLen = nameLen
			ensemble.InstrumentVelocityData.Base.Name = name

			// Ensemble Data Set
			ensemble.InstrumentVelocityData.Decode(data[packetPointer : packetPointer+dataSetSize])

			// Move to the next dataset
			packetPointer += dataSetSize
		} else if strings.Contains(name, earthVelocityID) {
			// Base
			ensemble.EarthVelocityData.Base.Enstype = enstype
			ensemble.EarthVelocityData.Base.NumElements = numElements
			ensemble.EarthVelocityData.Base.ElementMultiplier = elementMultiplier
			ensemble.EarthVelocityData.Base.Imag = imag
			ensemble.EarthVelocityData.Base.NameLen = nameLen
			ensemble.EarthVelocityData.Base.Name = name

			// Ensemble Data Set
			ensemble.EarthVelocityData.Decode(data[packetPointer : packetPointer+dataSetSize])

			// Move to the next dataset
			packetPointer += dataSetSize
		} else if strings.Contains(name, amplitudeID) {
			// Base
			ensemble.AmplitudeData.Base.Enstype = enstype
			ensemble.AmplitudeData.Base.NumElements = numElements
			ensemble.AmplitudeData.Base.ElementMultiplier = elementMultiplier
			ensemble.AmplitudeData.Base.Imag = imag
			ensemble.AmplitudeData.Base.NameLen = nameLen
			ensemble.AmplitudeData.Base.Name = name

			// Ensemble Data Set
			ensemble.AmplitudeData.Decode(data[packetPointer : packetPointer+dataSetSize])

			// Move to the next dataset
			packetPointer += dataSetSize
		} else if strings.Contains(name, correlationID) {
			// Base
			ensemble.CorrelationData.Base.Enstype = enstype
			ensemble.CorrelationData.Base.NumElements = numElements
			ensemble.CorrelationData.Base.ElementMultiplier = elementMultiplier
			ensemble.CorrelationData.Base.Imag = imag
			ensemble.CorrelationData.Base.NameLen = nameLen
			ensemble.CorrelationData.Base.Name = name

			// Ensemble Data Set
			ensemble.CorrelationData.Decode(data[packetPointer : packetPointer+dataSetSize])

			// Move to the next dataset
			packetPointer += dataSetSize
		} else {
			// Move to the next dataset
			packetPointer += dataSetSize
		}

		// Check if there is anymore data
		if packetPointer+4 >= len(data) {
			break
		}

	}

	return ensemble

}

// GenerateIndex will find the location of the data within
// the slice based off the index value given.
func GenerateIndex(index int, nameLen uint32, ensType uint32) int {
	datatype := BytesInFloat

	switch ensType {
	case dataTypeByte:
		datatype = BytesInInt8
		break
	case dataTypeInt:
		datatype = BytesInInt32
		break
	case dataTypeFloat:
		datatype = BytesInFloat
		break
	default:
		datatype = BytesInFloat
		break
	}

	return getHeaderSize(nameLen) + (index * datatype)
}

// GetBinBeamIndex will get the index in the data for the value.
// This is used for data with bin data, like velocity, correlation and amplitde.
func GetBinBeamIndex(nameLen int, numBins int, beam int, bin int) int {
	return getHeaderSize(uint32(nameLen)) + (beam * numBins * BytesInFloat) + (bin * BytesInFloat)
}

// getDataSetSize will get the size of the dataset.  It will user the number of elements and the
// element mulitipler to determine how many bytes are within the dataset.  It will also include the
// header.  It will also need to know how many bytes per element.
// Data is usually numElements x elementMultipler
// Bin data: bins x beams
// Other data: numElements x 1
func getDataSetSize(ensType uint32, nameLen uint32, numElements uint32, elementMultiplier uint32) int {

	datatype := BytesInFloat

	switch ensType {
	case dataTypeByte:
		datatype = BytesInInt8
		break
	case dataTypeInt:
		datatype = BytesInInt32
		break
	case dataTypeFloat:
		datatype = BytesInFloat
		break
	default:
		datatype = BytesInFloat
		break
	}

	return ((int(numElements) * int(elementMultiplier)) * int(datatype)) + getHeaderSize(nameLen)
}

// getHeaderSize will get the number bytes in the header.
func getHeaderSize(nameLen uint32) int {
	return int(nameLen) + (BytesInInt32 * (numDataSetHeaderElements - 1))
}

// findHeader will find the header to the ensemble.
func findHeader() {
	//log.Printf("Search for ID %d Index: %d", buffer.Len(), headerIndex)

	// Remove the first byte until it is an ID
	for {
		// Read the next byte in the buffer
		id, err := buffer.ReadByte()

		// Check for errors or if the buffer is empty
		if err != nil {
			//log.Print("Error finding ID. " + err.Error())
			return
		}

		// Verify still in range
		if headerIndex >= len(header) {
			return
		}

		// Check if the ID is correct
		if id != 0x80 {
			// Start over since not an ID value
			headerIndex = 0
		} else {
			//log.Printf("Header Index: %d", headerIndex)
			header[headerIndex] = id // Set the first ID
			headerIndex++

			// Found a complete ID
			if headerIndex == idlen {
				break // Break out for loop
			}
		}
	}

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
// It will reset all the values.
func clearHeader() {
	// for i := 0; i < 32; i++ {
	// 	header[i] = 0
	// }
	headerIndex = 0
	headerFound = false

	//time.Sleep(500 * time.Millisecond)

}

// calculateEnsembleChecksum will calculate the checksum for
// the given ensemble. This will use CRC-16 to calculate the checksum.
// Give all bytes in the Ensemble including the checksum.
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
