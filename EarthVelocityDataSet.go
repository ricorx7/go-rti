package rti

import (
	"encoding/binary"
	"math"
)

// EarthVelocityDataSet will contain all the Earth Velocity Data set values.
// These values describe water profile data in m/s.
// The data will be stored in array.  The array size will be based off the
// base data set.
// Bin x Beam
type EarthVelocityDataSet struct {
	Base     BaseDataSet      // Base Dataset
	Velocity [][]float32      // Velcity data in m/s
	Vectors  []VelocityVector // Velocity vector with maginitude and direction
}

// Decode will take the binary data and decode into
// into the ensemble data set.
func (vel *EarthVelocityDataSet) Decode(data []byte) {

	// Initialize the 2D array
	// [Bins][Beams]
	vel.Velocity = make([][]float32, vel.Base.NumElements)
	for i := range vel.Velocity {
		vel.Velocity[i] = make([]float32, vel.Base.ElementMultiplier)
	}
	vel.Vectors = make([]VelocityVector, vel.Base.NumElements)

	// Not enough data
	if uint32(len(data)) < vel.Base.NumElements*vel.Base.ElementMultiplier*uint32(BytesInFloat) {
		return
	}

	// Set each beam and bin data
	ptr := 0
	for beam := 0; beam < int(vel.Base.ElementMultiplier); beam++ {
		for bin := 0; bin < int(vel.Base.NumElements); bin++ {
			// Get the location of the data
			ptr = GetBinBeamIndex(int(vel.Base.NameLen), int(vel.Base.NumElements), beam, bin)

			// Set the data to float
			bits := binary.LittleEndian.Uint32(data[ptr : ptr+BytesInFloat])
			vel.Velocity[bin][beam] = math.Float32frombits(bits)

			// Calculate the velocity vector
			vel.Vectors[bin] = calcVV(vel.Velocity[bin])
		}
	}
}

// calcVV will calculate the velocity vector for each bin.
func calcVV(binData []float32) VelocityVector {
	// Init the values
	var vv = VelocityVector{Magnitude: BadVelocity, DirectionXNorth: BadVelocity, DirectionYNorth: BadVelocity}

	// Need East, North, Vert
	if len(binData) < 3 {
		return vv
	}

	// All the data must be good
	if binData[0] == BadVelocity || binData[1] == BadVelocity || binData[2] == BadVelocity {
		return vv
	}

	// Magnitude
	// Sqrt(East^2 + North^2 + Vert^2)
	mag := math.Sqrt(math.Pow(float64(binData[0]), 2) + math.Pow(float64(binData[1]), 2) + math.Pow(float64(binData[2]), 2))
	vv.Magnitude = math.Abs(mag)

	// Direct X North
	// atan2(east, north)
	dirXNorth := math.Atan2(float64(binData[0]), float64(binData[1])) * (180.0 / math.Pi)
	if dirXNorth < 0.0 {
		dirXNorth = 360.0 + dirXNorth
	}
	vv.DirectionXNorth = dirXNorth

	// Direct Y North
	// atan2(north, east)
	dirYNorth := math.Atan2(float64(binData[1]), float64(binData[2])) * (180.0 / math.Pi)
	if dirYNorth < 0.0 {
		dirYNorth = 360.0 + dirYNorth
	}
	vv.DirectionYNorth = dirYNorth

	return vv
}
