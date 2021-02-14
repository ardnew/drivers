package main

import (
	"errors"
	"strconv"
)

const (
	dimensionMax = Dimension(32767)
	dimensionMin = Dimension(0)
)

var ErrInvalidDimension = errors.New(
	"invalid dimension (" +
		strconv.FormatInt(int64(dimensionMin), 10) + "-" +
		strconv.FormatInt(int64(dimensionMax), 10) + ")",
)

// Dimension represents one component of a discrete Pixel's 2D coordinates
type Dimension uint16 // (physical space)

// Position converts the receiver d in physical space to a component of a
// Particle coordinate in logical space.
func (d Dimension) Position() Position {
	return Position(int(d) * int(velocityMax))
}

// Pixel represents an object with fixed coordinates in space.
//
// See the godoc comment on type Particle for details about the two coordinate
// spaces used to describe objects.
type Pixel struct {
	x, y Dimension
}
