package julia

import (
	"math"
	"testing"
)

func TestIterate(t *testing.T) {
	tests := []struct {
		name         string
		z0, c        complex128
		maxIter      int
		escapeRadius float64
		wantEscaped  bool
	}{
		{
			name:         "origin with zero c does not escape",
			z0:           0,
			c:            0,
			maxIter:      256,
			escapeRadius: DefaultEscapeRadius,
			wantEscaped:  false,
		},
		{
			name:         "large z escapes immediately",
			z0:           10 + 0i,
			c:            0,
			maxIter:      256,
			escapeRadius: DefaultEscapeRadius,
			wantEscaped:  true,
		},
		{
			name:         "classic Julia c=-1 cycle from origin does not escape",
			z0:           0,
			c:            -1,
			maxIter:      256,
			escapeRadius: DefaultEscapeRadius,
			wantEscaped:  false,
		},
		{
			name:         "just outside escape radius escapes",
			z0:           2.01 + 0i,
			c:            0,
			maxIter:      256,
			escapeRadius: DefaultEscapeRadius,
			wantEscaped:  true,
		},
		{
			name:         "point inside unit disk with c=0 does not escape",
			z0:           0.5 + 0i,
			c:            0,
			maxIter:      256,
			escapeRadius: DefaultEscapeRadius,
			wantEscaped:  false,
		},
		{
			name:         "escaped point has non-negative smooth value",
			z0:           0.5 + 0.5i,
			c:            0.4 + 0.4i,
			maxIter:      256,
			escapeRadius: DefaultEscapeRadius,
			wantEscaped:  true,
		},
		{
			name:         "small c near origin does not escape",
			z0:           0,
			c:            -0.1 + 0.1i,
			maxIter:      256,
			escapeRadius: DefaultEscapeRadius,
			wantEscaped:  false,
		},
		{
			name:         "boundary just inside escape radius with c=0 escapes after iterations",
			z0:           1.99 + 0i,
			c:            0,
			maxIter:      256,
			escapeRadius: DefaultEscapeRadius,
			wantEscaped:  true,
		},
		{
			name:         "overflow to NaN is treated as escaped not interior",
			z0:           1e155 + 1e155i,
			c:            0,
			maxIter:      256,
			escapeRadius: DefaultEscapeRadius,
			wantEscaped:  true,
		},
		{
			name:         "small escape radius with |z|<1 escape gives finite smooth",
			z0:           0.6 + 0i,
			c:            0,
			maxIter:      256,
			escapeRadius: 0.5, // er2=0.25; z0 escapes immediately with mag2=0.36, |z|<1
			wantEscaped:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			escaped, smooth := Iterate(tt.z0, tt.c, tt.maxIter, tt.escapeRadius)

			if escaped != tt.wantEscaped {
				t.Errorf("escaped = %v, want %v", escaped, tt.wantEscaped)
			}

			if escaped {
				if smooth < 0 || math.IsNaN(smooth) || math.IsInf(smooth, 0) {
					t.Errorf("smooth = %v, want finite non-negative value", smooth)
				}
			} else {
				if smooth != -1.0 {
					t.Errorf("smooth = %f, want -1.0 (interior sentinel)", smooth)
				}
			}
		})
	}
}

func TestPixelToComplex(t *testing.T) {
	p := Params{
		MinX: -2,
		MaxX: 2,
		MinY: -1.5,
		MaxY: 1.5,
	}

	tests := []struct {
		name   string
		px, py int
		width  int
		height int
		wantRe float64
		wantIm float64
	}{
		{
			name:   "top-left corner maps to (minX, minY)",
			px:     0,
			py:     0,
			width:  100,
			height: 100,
			wantRe: -2.0,
			wantIm: -1.5,
		},
		{
			name:   "bottom-right corner maps near (maxX, maxY)",
			px:     99,
			py:     99,
			width:  100,
			height: 100,
			wantRe: -2.0 + 4.0*99.0/100.0,
			wantIm: -1.5 + 3.0*99.0/100.0,
		},
		{
			name:   "center pixel maps to center",
			px:     50,
			py:     50,
			width:  100,
			height: 100,
			wantRe: 0.0,
			wantIm: 0.0,
		},
		{
			name:   "1x1 image maps to top-left",
			px:     0,
			py:     0,
			width:  1,
			height: 1,
			wantRe: -2.0,
			wantIm: -1.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := PixelToComplex(tt.px, tt.py, tt.width, tt.height, p)
			re := real(z)
			im := imag(z)
			const eps = 1e-10

			if math.Abs(re-tt.wantRe) > eps {
				t.Errorf("real = %f, want %f", re, tt.wantRe)
			}
			if math.Abs(im-tt.wantIm) > eps {
				t.Errorf("imag = %f, want %f", im, tt.wantIm)
			}
		})
	}
}
