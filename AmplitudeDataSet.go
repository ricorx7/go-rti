package rti

import (
	"encoding/binary"
	"math"
)

// AmplitudeDataSet will contain all the Amplitude Data set values.
// These values describe water profile signal strength in dB.
// The data will be stored in array.  The array size will be based off the
// base data set.
// Bin x Beam
type AmplitudeDataSet struct {
	Base      BaseDataSet // Base Dataset
	Amplitude [][]float32 // Amplitude data in dB
}

// Decode will take the binary data and decode into
// into the ensemble data set.
func (amp *AmplitudeDataSet) Decode(data []byte) {

	// Initialize the 2D array
	// [Bins][Beams]
	amp.Amplitude = make([][]float32, amp.Base.NumElements)
	for i := range amp.Amplitude {
		amp.Amplitude[i] = make([]float32, amp.Base.ElementMultiplier)
	}

	// Not enough data
	if uint32(len(data)) < amp.Base.NumElements*amp.Base.ElementMultiplier*uint32(BytesInFloat) {
		return
	}

	// Set each beam and bin data
	ptr := 0
	for beam := 0; beam < int(amp.Base.ElementMultiplier); beam++ {
		for bin := 0; bin < int(amp.Base.NumElements); bin++ {
			// Get the location of the data
			ptr = GetBinBeamIndex(int(amp.Base.NameLen), int(amp.Base.NumElements), beam, bin)

			// Set the data to float
			bits := binary.LittleEndian.Uint32(data[ptr : ptr+BytesInFloat])
			amp.Amplitude[bin][beam] = math.Float32frombits(bits)
		}
	}

}
