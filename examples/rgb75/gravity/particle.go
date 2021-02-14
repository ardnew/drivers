package main

import (
	"math"
	"math/rand"
)

const (
	dimensionMax = Dimension(32767)
	dimensionMin = Dimension(0)
	positionMax  = Position(int32(velocityMax) * int32(dimensionMax))
	positionMin  = Position(0)
	velocityMax  = Velocity(256)
	velocityMin  = -velocityMax
	velocityMax2 = velocityMax * velocityMax // max²
)

// Dimension represents one component of a discrete pixel's 2D coordinates
type Dimension uint16 // (physical space)

// Position converts the receiver d in physical space to a component of a
// Particle coordinate in logical space.
func (d Dimension) Position() Position {
	return Position(int(d) * int(velocityMax))
}

// Position represents one component of a discrete Particle's 2D coordinates
type Position int32 // (logical space)

// Dimension converts the receiver p in logical space to a component of a real
// pixel coordinate in physical space.
func (p Position) Dimension() Dimension {
	return Dimension(int(p) / int(velocityMax))
}

// Move returns the receiver Position p, and its equivalent Dimension, adjusted
// by given Velocity v.
func (p Position) Move(v Velocity) (Position, Dimension) {
	pos := Position(int(p) + int(v))
	return pos, pos.Dimension()
}

// Velocity represents one component of a discrete Particle's 2D velocity
type Velocity int32 // (logical space)

// Reverse returns the receiver Velocity v in the opposite direction and scaled
// by a given elasticity.
func (v Velocity) Reverse(elasticity int) Velocity {
	return Velocity(int(-v) * elasticity / int(velocityMax))
}

// Abs returns the absolute value of the receiver Velocity v.
func (v Velocity) Abs() Velocity {
	if v < 0 {
		return -v
	}
	return v
}

// Particle represents an object moving through space.
//
// The space through which a Particle moves is referred to in documentation as
// "logical space", since that space is much larger than the "physical space"
// used to describe physical pixel coordinates; these added logical coordinates
// exist "in-between" pixels, and allow for smoother movement in the absence of
// floating-point coordinates.
//
// Particles in logical space are always eventually projected onto physical
// space when displaying them with a pixel.
type Particle struct {
	ix, iy Dimension
	px, py Position
	vx, vy Velocity
}

// Particles represents all Particle objects in the Field's 2D space.
type Particles []Particle

// ParticleMove defines a callback used to notify callers when a Particle moves.
type ParticleMove func(f *Field, p *Particle, x, y Dimension)

// MakeParticles returns a new Particle buffer of given Field f and count n.
// Each Particle is initially positioned in the first unoccupied pixel on the
// Field.
func MakeParticles(f *Field, n int) Particles {
	if n >= 0 {
		particle := make(Particles, n)
		for i := range particle {
			x := Dimension(i) % f.width
			y := Dimension(i) / f.width
			particle[i].SetPosition(f, x, y, x.Position(), y.Position())
		}
		return particle
	}
	return nil
}

// Accelerate applies the current acceleration due to gravity to the velocity of
// receiver p, with a slight perturbation epsilon.
// This only changes the Particle velocity; it does not affect its Position.
func (p *Particle) Accelerate(x, y, z, epsilon int) {
	// apply random perturbation to the values read from accelerometer.
	// do not use MakeVelocity, as it will prematurely clip the x, y components.
	p.vx += Velocity(x + rand.Intn(epsilon))
	p.vy += Velocity(y + rand.Intn(epsilon))
	// clip the resulting vector to maximum velocity
	v2 := p.vx*p.vx + p.vy*p.vy
	if v2 > velocityMax2 { // implies v > velocityMax in some direction
		v := math.Sqrt(float64(v2))
		p.vx = Velocity(int(float64(velocityMax*p.vx) / v))
		p.vy = Velocity(int(float64(velocityMax*p.vy) / v))
	}
}

// SetPosition sets the coordinates of the receiver Particle p, and updates the
// Obstacle coordinates of the given Field f.
func (p *Particle) SetPosition(f *Field, ix, iy Dimension, px, py Position) {
	if nil != f.handleMove {
		f.handleMove(f, p, ix, iy)
	}
	f.obstacle.Clr(p.ix, p.iy)
	p.ix, p.iy = ix, iy
	p.px, p.py = px, py
	f.obstacle.Set(p.ix, p.iy)
}

// Move attempts to change the coordinates of the receiver Particle p based on
// its current velocity, or reverses velocity if the change would collide with
// an Obstacle.
func (p *Particle) Move(f *Field) {

	// first, compute destination Position based on current Velocity
	px, ix := p.px.Move(p.vx)
	py, iy := p.py.Move(p.vy)

	// next, verify we are moving within Field boundaries
	if px < 0 {
		p.vx = p.vx.Reverse(f.elasticity)
		px, ix = 0, 0
	} else if px > f.xMax {
		p.vx = p.vx.Reverse(f.elasticity)
		px, ix = f.xMax, f.width-1
	}
	if py < 0 {
		p.vy = p.vy.Reverse(f.elasticity)
		py, iy = 0, 0
	} else if py > f.yMax {
		p.vy = p.vy.Reverse(f.elasticity)
		py, iy = f.yMax, f.height-1
	}

	// then, determine if we are moving into a new real pixel in physical space
	if dp := f.PixelIndex(p.ix, p.iy) - f.PixelIndex(ix, iy); 0 != dp {
		// check if the destination pixel contains an Obstacle
		if f.obstacle.Get(ix, iy) {
			if dp < 0 {
				dp = -dp // absolute value of index difference
			}
			// determine which direction the Obstacle exists
			switch dp {
			case 1: // obstructed by 1 pixel to the left or right
				p.vx = p.vx.Reverse(f.elasticity)
				px, ix = p.px, p.ix

			case int(f.width): // obstructed by 1 pixel to the top or bottom
				p.vy = p.vy.Reverse(f.elasticity)
				py, iy = p.py, p.iy

			default: // obstructed by 1 pixel in a diagonal direction
				if p.vx.Abs() >= p.vy.Abs() {
					if !f.obstacle.Get(ix, p.iy) {
						p.vy = p.vy.Reverse(f.elasticity)
						py, iy = p.py, p.iy
					} else if !f.obstacle.Get(p.ix, iy) {
						p.vx = p.vx.Reverse(f.elasticity)
						px, ix = p.px, p.ix
					} else {
						p.vx = p.vx.Reverse(f.elasticity)
						p.vy = p.vy.Reverse(f.elasticity)
						px, ix = p.px, p.ix
						py, iy = p.py, p.iy
					}
				} else {
					if !f.obstacle.Get(p.ix, iy) {
						p.vx = p.vx.Reverse(f.elasticity)
						px, ix = p.px, p.ix
					} else if !f.obstacle.Get(ix, p.iy) {
						p.vy = p.vy.Reverse(f.elasticity)
						py, iy = p.py, p.iy
					} else {
						p.vx = p.vx.Reverse(f.elasticity)
						p.vy = p.vy.Reverse(f.elasticity)
						px, ix = p.px, p.ix
						py, iy = p.py, p.iy
					}
				}
			}
		}
	}

	// update coordinates of both Particle and Obstacle.
	p.SetPosition(f, ix, iy, px, py)
}
