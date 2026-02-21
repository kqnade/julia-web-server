package renderer

import (
	"testing"

	"github.com/kqnade/julia-web-server/internal/julia"
)

func defaultParams(width, height int) julia.Params {
	return julia.Params{
		MinX:         -2,
		MaxX:         2,
		MinY:         -1.5,
		MaxY:         1.5,
		C:            -0.7 + 0.27015i,
		Width:        width,
		Height:       height,
		MaxIter:      julia.DefaultMaxIter,
		EscapeRadius: julia.DefaultEscapeRadius,
	}
}

func TestRender_BufferSize(t *testing.T) {
	p := defaultParams(100, 80)
	buf := Render(p)
	want := 100 * 80
	if len(buf) != want {
		t.Errorf("buffer length = %d, want %d", len(buf), want)
	}
}

func TestRender_AllEscaping(t *testing.T) {
	// Region far from origin: all points should escape
	p := julia.Params{
		MinX:         10,
		MaxX:         12,
		MinY:         10,
		MaxY:         12,
		C:            0,
		Width:        16,
		Height:       16,
		MaxIter:      256,
		EscapeRadius: julia.DefaultEscapeRadius,
	}
	buf := Render(p)
	for i, v := range buf {
		if v < 0 {
			t.Errorf("buf[%d] = %f, want >= 0 (all points should escape)", i, v)
			break
		}
	}
}

func TestRender_InteriorPoints(t *testing.T) {
	// Small region near origin with c=0: points with |z|<1 don't escape
	p := julia.Params{
		MinX:         -0.1,
		MaxX:         0.1,
		MinY:         -0.1,
		MaxY:         0.1,
		C:            0,
		Width:        8,
		Height:       8,
		MaxIter:      256,
		EscapeRadius: julia.DefaultEscapeRadius,
	}
	buf := Render(p)
	for i, v := range buf {
		if v != -1.0 {
			t.Errorf("buf[%d] = %f, want -1.0 (interior points)", i, v)
			break
		}
	}
}

func TestRender_Deterministic(t *testing.T) {
	p := defaultParams(64, 64)
	buf1 := Render(p)
	buf2 := Render(p)

	if len(buf1) != len(buf2) {
		t.Fatalf("different lengths: %d vs %d", len(buf1), len(buf2))
	}
	for i := range buf1 {
		if buf1[i] != buf2[i] {
			t.Errorf("buf[%d] differs: %f vs %f", i, buf1[i], buf2[i])
			break
		}
	}
}

func TestRender_1x1(t *testing.T) {
	p := defaultParams(1, 1)
	buf := Render(p)
	if len(buf) != 1 {
		t.Errorf("buffer length = %d, want 1", len(buf))
	}
}

func TestRender_NonSquare(t *testing.T) {
	p := defaultParams(200, 50)
	buf := Render(p)
	want := 200 * 50
	if len(buf) != want {
		t.Errorf("buffer length = %d, want %d", len(buf), want)
	}
}

func TestRender_ZeroHeight_NoPanic(t *testing.T) {
	p := defaultParams(64, 0)
	buf := Render(p)
	if len(buf) != 0 {
		t.Errorf("buffer length = %d, want 0", len(buf))
	}
}
