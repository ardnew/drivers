// Package conf provides the interface used to modify INA260 configuration.
package conf // import "tinygo.org/x/drivers/ina260/conf"

// SampleSize determines the number of samples that are collected and averaged.
// Each sample requires a given amount of ADC conversion time, and there are
// only a small number of discrete times that may be selected (ConversionTime).
type SampleSize uint8

// Constants representing each possible value of type SampleSize.
const (
	SampleSize1       SampleSize = 0 // (000b) -- default
	SampleSize4       SampleSize = 1 // (001b)
	SampleSize16      SampleSize = 2 // (010b)
	SampleSize64      SampleSize = 3 // (011b)
	SampleSize128     SampleSize = 4 // (100b)
	SampleSize256     SampleSize = 5 // (101b)
	SampleSize512     SampleSize = 6 // (110b)
	SampleSize1024    SampleSize = 7 // (111b)
	SampleSizeDefault SampleSize = SampleSize1
)

// ConversionTime sets the conversion time for the voltage and current
// measurement. These are the only intervals recognized by the INA260 hardware
// (per the datasheet). One complete lapse in the selected duration represents 1
// sample. Therefore, the total time required for a single measurement is
// calculated as selected conversion time (ConversionTime) multiplied by the
// selected number of samples (SampleSize).
type ConversionTime uint8

// Constants representing each possible value of type ConversionTime.
const (
	ConversionTime140us   ConversionTime = 0 // (000b)
	ConversionTime204us   ConversionTime = 1 // (001b)
	ConversionTime332us   ConversionTime = 2 // (010b)
	ConversionTime588us   ConversionTime = 3 // (011b)
	ConversionTime1p1ms   ConversionTime = 4 // (100b) -- default (voltage, current)
	ConversionTime2p116ms ConversionTime = 5 // (101b)
	ConversionTime4p156ms ConversionTime = 6 // (110b)
	ConversionTime8p244ms ConversionTime = 7 // (111b)
	ConversionTimeDefault ConversionTime = ConversionTime1p1ms
)

// OperatingMode determines how measurements should be performed and updated in
// internal data registers. You can read the contents at any time in both modes,
// but this will affect when they update.
type OperatingMode uint8

// Constants representing each possible value of type OperatingMode.
const (
	OperatingModeTriggered  OperatingMode = 0 // (000b)
	OperatingModeContinuous OperatingMode = 1 // (001b) -- default
	OperatingModeDefault    OperatingMode = OperatingModeContinuous
)

// OperatingType specifies which measurements are performed for each conversion.
// You can perform current-only, voltage-only, or both current and voltage. Note
// the bit patterns result in the following equalities:
//   OperatingTypeShutdown == 0
//   OperatingTypePower    == ( OperatingTypeVoltage | OperatingTypeCurrent )
type OperatingType uint8

// Constants representing each possible value of type OperatingType.
const (
	OperatingTypeShutdown OperatingType = 0 // (000b)
	OperatingTypeCurrent  OperatingType = 1 // (001b)
	OperatingTypeVoltage  OperatingType = 2 // (010b)
	OperatingTypePower    OperatingType = 3 // (011b) -- default
	OperatingTypeDefault  OperatingType = OperatingTypePower
)

// Configuration represents content of the CONFIGURATION register (00h).
type Configuration struct {
	Type  OperatingType
	Mode  OperatingMode
	Ctime ConversionTime
	Vtime ConversionTime
	Size  SampleSize
}

// Default stores the default configuration of an INA260.
func Default() Configuration {
	return Configuration{
		Type:  OperatingTypeDefault,
		Mode:  OperatingModeDefault,
		Ctime: ConversionTimeDefault,
		Vtime: ConversionTimeDefault,
		Size:  SampleSizeDefault,
	}
}
