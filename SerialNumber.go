package rti

// SerialNumber is the ADCP serial number.
// It includes the hardware type, Subsystem Configurations
// Spare and the serial number.
type SerialNumber struct {
	Hardware   AdcpBaseHardware // Hardware type
	Spare      string           // Spare String
	SerialNum  int              // Serial Number
	Subsystems []Subsystem      // Subystems
}

// AdcpBaseHardware is the ADCP Base Hardware Enum Type
type AdcpBaseHardware int

// ADCP Base Hardware Enum value
const (
	AdcpHardware00 AdcpBaseHardware = 1 + iota // ADCP Base Hardware 00
	AdcpHardware01                  = iota     // ADCP Base Hardware 01
	AdcpHardware02                  = iota     // ADCP Base Hardware 02
)

// adcpBaseHardwares is the string representation of the Base Hardware enum
var adcpBaseHardwares = [...]string{
	"00",
	"01",
	"02",
}

// ToString is the ADCP Hardware to string
func (h AdcpBaseHardware) ToString() string { return adcpBaseHardwares[h-1] }

// Subsystem is the ADCP Subsystem Configuration
type Subsystem byte

const (
	// SubEmpty -  Empty ASCI: 0 HEX: 0x0
	SubEmpty Subsystem = 0x0
	// SubSpare0 - Spare ASCII: 0  HEX: 30
	SubSpare0 Subsystem = 0x30
	// Sub2Mhz4Beam20DegPiston1 -  2MHz 4Beam 20 Degree Piston  ASCII: 1 HEX: 0x31
	Sub2Mhz4Beam20DegPiston1 Subsystem = 0x31
)

// ToString is the ADCP Subystem to string
func (ss Subsystem) ToString() string {
	switch ss {
	case SubEmpty:
		return "Empty"
	case SubSpare0:
		return "Spare"
	case Sub2Mhz4Beam20DegPiston1:
		return "2MHz 4Beam 20 Degree Piston"
	default:
		return ""
	}
}
