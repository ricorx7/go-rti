package rti

// SubsystemConfiguration is the ADCP subsystem configuration.
// An  ADCP can contain multiple frequencies and configurations.  A subsystem is the
// frequency type.  The configuration is the configuration of that
// frequency.
type SubsystemConfiguration struct {
	CepoIndex uint8 // Index within the CEPO command
}

// Decode will decode the binary data given into the Subsystem configuration.
func (ssConfig *SubsystemConfiguration) Decode(data []byte) {
	// Ensure enough data
	if len(data) < 4 {
		return
	}

	// CEPO index
	ssConfig.CepoIndex = uint8(data[3])

}
