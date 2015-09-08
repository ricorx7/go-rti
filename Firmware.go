package rti

// Firmware is the ADCP version.
// It includes the subsystem configuraiton.
type Firmware struct {
	Major         uint8 // Firmware Major number
	Minior        uint8 // Firmware Minor number
	Revision      uint8 // Firmware Revision
	SubsystemCode byte  // Subsystem code
}

// Decode will decode the binary data given into the firmware version
// and subsystem code.
func (firm *Firmware) Decode(data []byte) {
	// Ensure enough data
	if len(data) < 4 {
		return
	}

	firm.Major = uint8(data[0])
	firm.Minior = uint8(data[1])
	firm.Revision = uint8(data[2])
	firm.SubsystemCode = data[3]
}
