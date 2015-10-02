package rti

// MaxNumDataSets is the number number of datasets.
const MaxNumDataSets = 12

// BytesInInt32 is the number of bytes in Int32
const BytesInInt32 = 4

// BytesInFloat is the number of bytes in float
const BytesInFloat = 4

// BytesInInt8 is the number of bytes in Int8
const BytesInInt8 = 1

// BAD_VELOCITY is the value for a bad value.
const BadVelocity = 88.88

// VelocityVector holds the water maginitude and direction.
type VelocityVector struct {
	Magnitude       float64 // Magnitude of the water velocity
	DirectionXNorth float64 // Direction of the water velocity with X north
	DirectionYNorth float64 // Direction of the water velocity with Y north
}

// BaseDataSet is the base for all data sets.
// It contains all the values to know the type
// and size of the data set.
type BaseDataSet struct {
	Enstype           uint32
	NumElements       uint32
	ElementMultiplier uint32
	Imag              uint32
	NameLen           uint32
	Name              string
}

// Ensemble is the container for the ADCP data.
// It will contain all the datasets.
type Ensemble struct {
	binaryData []byte // Binary Data

	EnsembleData           EnsembleDataSet           // Ensemble Data Set
	AncillaryData          AncillaryDataSet          // Ancillary Data Set
	BeamVelocityData       BeamVelocityDataSet       // Beam Velocity Data Set
	InstrumentVelocityData InstrumentVelocityDataSet // Instrument Velocity Data Set
	EarthVelocityData      EarthVelocityDataSet      // Earth Velocity Data Set
	AmplitudeData          AmplitudeDataSet          // Amplitude Data Set
	CorrelationData        CorrelationDataSet        // Correlation Data Set
}
