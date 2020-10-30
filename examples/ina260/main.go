package main

import (
	"machine"

	"tinygo.org/x/drivers/ina260"
)

func init() {
	machine.I2C1.Configure(machine.I2CConfig{
		Frequency: machine.TWI_FREQ_400KHZ,
		SDA:       machine.I2C1_SDA_PIN,
		SCL:       machine.I2C1_SCL_PIN,
	})
}

// BUG
// declaring this var here causes the following error at compilation:
//
//   # tinygo.org/x/drivers/examples/ina260
//   tinygo.org/x/drivers/examples/ina260/<init>: interp: unknown GEP
//
//   traceback:
//   tinygo.org/x/drivers/examples/ina260/<init>:
//     %3 = getelementptr inbounds { %machine.I2C }, { %machine.I2C }* %2, i32 0, i32 0, !dbg !4298
//
// (see other declaration below in `main`)
var pow = ina260.New(machine.I2C1)

func main() {

	// BUG
	// declaring the device here compiles OK
	//var pow = ina260.New(machine.I2C1)

	println("initializing INA260 on I2C1")

	if pow.Connected() {
		println("connected: INA260")
	} else {
		println("INA260 not connected!")
	}
}
