package rti

// EnsembleDataSet will contain all the Ensemble Data set values.
// These values describe ensemble with integer values.  This includes
// the serial number and firmware version.
type EnsembleDataSet struct {
	EnsembleNumber   int          // Ensemble number
	NumBins          int          // Number of bins
	NumBeams         int          // Number of beams
	DesiredPingCount int          // Desired Number of pings
	ActualPingCount  int          // Actual Number of pings
	SerialNumber     SerialNumber // Serial Number
}
