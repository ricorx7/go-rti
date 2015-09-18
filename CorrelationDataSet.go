package rti

import (
	"encoding/binary"
	"math"
)

// CorrelationDataSet will contain all the Correlation Data set values.
// These values describe water profile signal quality in percent (%).
// The data will be stored in array.  The array size will be based off the
// base data set.
// Bin x Beam
type CorrelationDataSet struct {
	Base        BaseDataSet // Base Dataset
	Correlation [][]float32 // Correlation data in %
}

// Decode will take the binary data and decode into
// into the ensemble data set.
func (corr *CorrelationDataSet) Decode(data []byte) {

	// Initialize the 2D array
	// [Bins][Beams]
	corr.Correlation = make([][]float32, corr.Base.NumElements)
	for i := range corr.Correlation {
		corr.Correlation[i] = make([]float32, corr.Base.ElementMultiplier)
	}

	// Not enough data
	if uint32(len(data)) < corr.Base.NumElements*corr.Base.ElementMultiplier*uint32(BytesInFloat) {
		return
	}

	// Set each beam and bin data
	ptr := 0
	for beam := 0; beam < int(corr.Base.ElementMultiplier); beam++ {
		for bin := 0; bin < int(corr.Base.NumElements); bin++ {
			// Get the location of the data
			ptr = GetBinBeamIndex(int(corr.Base.NameLen), int(corr.Base.NumElements), beam, bin)

			// Set the data to float
			bits := binary.LittleEndian.Uint32(data[ptr : ptr+BytesInFloat])
			corr.Correlation[bin][beam] = math.Float32frombits(bits)
		}
	}

}
