package ina260

// Constants defining the important register addresses of an INA260.
const (
	RegisterConfig  uint8 = 0x00
	RegisterCurrent uint8 = 0x01
	RegisterVoltage uint8 = 0x02
	RegisterPower   uint8 = 0x03
	RegisterMaskEn  uint8 = 0x06
	RegisterAlrtlim uint8 = 0x07
	RegisterMfgID   uint8 = 0xFE
	RegisterDevID   uint8 = 0xFF
)
