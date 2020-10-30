// Package ina260 provides a driver for the INA260 voltage/current sensor by TI.
//
// Datasheet: https://www.ti.com/lit/gpn/ina260
package ina260 // import "tinygo.org/x/drivers/ina260"

import (
	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/ina260/conf"
)

// Address is the I2C slave address of the INA260 used when creating a new
// connection.
// Use the GetAddress method on a Device receiver to obtain the slave address of
// a connected INA260.
var Address uint16 = 0x40

// Device wraps an I2C connection to an INA260 device.
type Device struct {
	bus     drivers.I2C
	address uint16
	config  conf.Configuration
}

// New creates a new INA260 connection. The I2C bus must already be configured.
func New(bus drivers.I2C) Device {
	return Device{
		bus:     bus,
		address: Address,
		config:  conf.Default(),
	}
}

// GetAddress returns the I2C slave address of the receiver d.
func (d *Device) GetAddress() uint16 {
	return d.address
}

// Connected returns whether we are communicating with the receiver d.
func (d *Device) Connected() bool {
	data := make([]byte, 2)
	d.bus.ReadRegister(uint8(d.address), RegisterDevID, data)
	id := uint16(data[0])<<8 | uint16(data[1])
	return id == 0x2270
}

// Configure modifies the configuration settings of the receiver d.
func (d *Device) Configure(config conf.Configuration) {
	d.config = config
}
