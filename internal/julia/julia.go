package julia

import "math"

const (
	DefaultMaxIter      = 256
	DefaultEscapeRadius = 2.0
)

// Params holds parameters for Julia set computation.
type Params struct {
	MinX, MaxX   float64
	MinY, MaxY   float64
	C            complex128
	Width        int
	Height       int
	MaxIter      int
	EscapeRadius float64
}

// Iterate performs the Julia set iteration starting from z0 with constant c.
// It returns whether the point escaped and the smooth iteration count.
// For escaped points, smooth >= 0 (clamped). For non-escaped points, smooth is -1.0.
func Iterate(z0, c complex128, maxIter int, escapeRadius float64) (escaped bool, smooth float64) {
	z := z0
	er2 := escapeRadius * escapeRadius

	for i := 0; i < maxIter; i++ {
		zr := real(z)
		zi := imag(z)
		mag2 := zr*zr + zi*zi

		// !(mag2 <= er2) catches both mag2 > er2 and NaN (from Inf-Inf overflow)
		if !(mag2 <= er2) {
			if math.IsNaN(mag2) || math.IsInf(mag2, 0) {
				return true, 0
			}
			// Smooth coloring: iteration + 1 - log(log(|z|)) / log(2)
			logMag := math.Log(mag2) / 2.0 // log(|z|) = log(mag2)/2
			// logMag must be > 0 (i.e. |z| > 1) for the formula to be valid.
			// With escapeRadius < 1, a point can escape with |z| <= 1, making
			// logMag <= 0 and math.Log(logMag) = NaN. Fall back to integer count.
			if logMag <= 0 {
				smooth = float64(i)
			} else {
				smooth = float64(i) + 1.0 - math.Log(logMag)/math.Log(2.0)
				if smooth < 0 {
					smooth = 0
				}
			}
			return true, smooth
		}

		z = z*z + c
	}

	return false, -1.0
}

// PixelToComplex converts pixel coordinates (px, py) to a complex number
// based on the given parameters. width and height must both be > 0.
func PixelToComplex(px, py, width, height int, p Params) complex128 {
	if width <= 0 || height <= 0 {
		panic("julia: PixelToComplex called with non-positive dimensions")
	}
	re := p.MinX + (p.MaxX-p.MinX)*float64(px)/float64(width)
	im := p.MinY + (p.MaxY-p.MinY)*float64(py)/float64(height)
	return complex(re, im)
}
