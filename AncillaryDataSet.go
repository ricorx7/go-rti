package rti

import (
	"encoding/binary"
	"math"
)

// AncillaryDataSet will contain all the Ancillary Data set values.
// These values describe ensemble with float values.
type AncillaryDataSet struct {
	Base            BaseDataSet // Base Dataset
	FirstBinRange   float32     // First bin location in meters
	BinSize         float32     // Bin size in meters
	FirstPingTime   float32     // First ping time in seconds
	LastPingTime    float32     // Last ping time in seconds
	Heading         float32     // Heading in degrees
	Pitch           float32     // Pitch in degrees
	Roll            float32     // Roll in degrees
	WaterTemp       float32     // Water temperature in degrees farenheit
	SystemTemp      float32     // System temperature in degrees farenheit
	Salinity        float32     // Salinity in Parts per Thousand (PPT)
	Pressure        float32     // Pressure in Pascals
	TransducerDepth float32     // Depth of the transducer in water in meters.  Used for speed of sound.
	SpeedOfSound    float32     // Speed of Sound in m/s
}

// Decode will take the binary data and decode into
// into the ensemble data set.
func (anc *AncillaryDataSet) Decode(data []byte) {

	// Not enough data
	if len(data) < 13*BytesInFloat {
		return
	}

	// First Bin Range
	ptr := GenerateIndex(0, anc.Base.NameLen, anc.Base.Enstype)
	bits := binary.LittleEndian.Uint32(data[ptr : ptr+BytesInFloat])
	anc.FirstBinRange = math.Float32frombits(bits)

	// Bin Size
	ptr = GenerateIndex(1, anc.Base.NameLen, anc.Base.Enstype)
	bits = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInFloat])
	anc.BinSize = math.Float32frombits(bits)

	// First Ping Time
	ptr = GenerateIndex(2, anc.Base.NameLen, anc.Base.Enstype)
	bits = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInFloat])
	anc.FirstPingTime = math.Float32frombits(bits)

	// Last Ping Time
	ptr = GenerateIndex(3, anc.Base.NameLen, anc.Base.Enstype)
	bits = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInFloat])
	anc.LastPingTime = math.Float32frombits(bits)

	// Heading
	ptr = GenerateIndex(4, anc.Base.NameLen, anc.Base.Enstype)
	bits = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInFloat])
	anc.Heading = math.Float32frombits(bits)

	// Pitch
	ptr = GenerateIndex(5, anc.Base.NameLen, anc.Base.Enstype)
	bits = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInFloat])
	anc.Pitch = math.Float32frombits(bits)

	// Roll
	ptr = GenerateIndex(6, anc.Base.NameLen, anc.Base.Enstype)
	bits = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInFloat])
	anc.Roll = math.Float32frombits(bits)

	// Water Temp
	ptr = GenerateIndex(7, anc.Base.NameLen, anc.Base.Enstype)
	bits = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInFloat])
	anc.WaterTemp = math.Float32frombits(bits)

	// System Temp
	ptr = GenerateIndex(8, anc.Base.NameLen, anc.Base.Enstype)
	bits = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInFloat])
	anc.SystemTemp = math.Float32frombits(bits)

	// Salinity
	ptr = GenerateIndex(9, anc.Base.NameLen, anc.Base.Enstype)
	bits = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInFloat])
	anc.Salinity = math.Float32frombits(bits)

	// Pressure
	ptr = GenerateIndex(10, anc.Base.NameLen, anc.Base.Enstype)
	bits = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInFloat])
	anc.Pressure = math.Float32frombits(bits)

	// Transducer Depth
	ptr = GenerateIndex(11, anc.Base.NameLen, anc.Base.Enstype)
	bits = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInFloat])
	anc.TransducerDepth = math.Float32frombits(bits)

	// Speed of Sound
	ptr = GenerateIndex(12, anc.Base.NameLen, anc.Base.Enstype)
	bits = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInFloat])
	anc.SpeedOfSound = math.Float32frombits(bits)

}
