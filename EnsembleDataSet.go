package rti

import "encoding/binary"

// EnsembleDataSet will contain all the Ensemble Data set values.
// These values describe ensemble with integer values.  This includes
// the serial number and firmware version.
type EnsembleDataSet struct {
	Base             BaseDataSet            // Base Dataset
	EnsembleNumber   uint32                 // Ensemble number
	NumBins          uint32                 // Number of bins
	NumBeams         uint32                 // Number of beams
	DesiredPingCount uint32                 // Desired Number of pings
	ActualPingCount  uint32                 // Actual Number of pings
	Status           uint32                 // ADCP Status
	Year             uint32                 // Year
	Month            uint32                 // Month
	Day              uint32                 // Day
	Hour             uint32                 // Hour
	Minute           uint32                 // Minute
	Second           uint32                 // Second
	HSec             uint32                 // Hundredth of second
	SerialNumber     SerialNumber           // Serial Number
	Firmware         Firmware               // Firmware
	SubsystemConfig  SubsystemConfiguration // Subsystem configuration
}

// Decode will take the binary data and decode into
// into the ensemble data set.
func (ens *EnsembleDataSet) Decode(data []byte) {

	// Not enough data
	if len(data) < 13*BytesInInt32 {
		return
	}

	// Ensemble number
	ptr := GenerateIndex(0, ens.Base.NameLen, ens.Base.Enstype)
	ens.EnsembleNumber = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInInt32])

	// Number of bins
	ptr = GenerateIndex(1, ens.Base.NameLen, ens.Base.Enstype)
	ens.NumBins = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInInt32])

	// Number of beams
	ptr = GenerateIndex(2, ens.Base.NameLen, ens.Base.Enstype)
	ens.NumBeams = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInInt32])

	// Desired Ping count
	ptr = GenerateIndex(3, ens.Base.NameLen, ens.Base.Enstype)
	ens.DesiredPingCount = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInInt32])

	// Actual Ping count
	ptr = GenerateIndex(4, ens.Base.NameLen, ens.Base.Enstype)
	ens.ActualPingCount = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInInt32])

	// Status
	ptr = GenerateIndex(5, ens.Base.NameLen, ens.Base.Enstype)
	ens.Status = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInInt32])

	// Year
	ptr = GenerateIndex(6, ens.Base.NameLen, ens.Base.Enstype)
	ens.Year = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInInt32])

	// Month
	ptr = GenerateIndex(7, ens.Base.NameLen, ens.Base.Enstype)
	ens.Month = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInInt32])

	// Day
	ptr = GenerateIndex(8, ens.Base.NameLen, ens.Base.Enstype)
	ens.Day = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInInt32])

	// Hour
	ptr = GenerateIndex(9, ens.Base.NameLen, ens.Base.Enstype)
	ens.Hour = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInInt32])

	// Minute
	ptr = GenerateIndex(10, ens.Base.NameLen, ens.Base.Enstype)
	ens.Minute = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInInt32])

	// Second
	ptr = GenerateIndex(11, ens.Base.NameLen, ens.Base.Enstype)
	ens.Second = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInInt32])

	// HSec
	ptr = GenerateIndex(12, ens.Base.NameLen, ens.Base.Enstype)
	ens.HSec = binary.LittleEndian.Uint32(data[ptr : ptr+BytesInInt32])

	// Ensure enough data is there to decode Firmware and serial number
	if len(data) >= 22*BytesInInt32 {
		// Serial Number
		ptr = GenerateIndex(13, ens.Base.NameLen, ens.Base.Enstype)
		ens.SerialNumber.SerialNumber = string(data[ptr : ptr+(8*BytesInInt32)])

		// Firmware
		ptr = GenerateIndex(21, ens.Base.NameLen, ens.Base.Enstype)
		ens.Firmware.Decode(data[ptr : ptr+BytesInInt32])
	}

	// Ensure enough data to decode Subsystem configuration
	if len(data) >= 23*BytesInInt32 {
		ptr = GenerateIndex(22, ens.Base.NameLen, ens.Base.Enstype)
		ens.SubsystemConfig.Decode(data[ptr : ptr+BytesInInt32])
	}
}
