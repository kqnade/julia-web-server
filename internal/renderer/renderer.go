package renderer

import (
	"runtime"
	"sync"

	"github.com/kqnade/julia-web-server/internal/julia"
)

// Render computes the Julia set for the given parameters and returns a
// float32 slice of length Width*Height in row-major order (left-to-right,
// top-to-bottom). Each value is the smooth iteration count (>= 0 for escaped
// points, -1.0 for interior points).
func Render(p julia.Params) []float32 {
	buf := make([]float32, p.Width*p.Height)

	numWorkers := runtime.NumCPU()
	if numWorkers > p.Height {
		numWorkers = p.Height
	}
	if numWorkers < 1 {
		numWorkers = 1
	}

	var wg sync.WaitGroup
	rowsPerWorker := p.Height / numWorkers

	for w := 0; w < numWorkers; w++ {
		startRow := w * rowsPerWorker
		endRow := startRow + rowsPerWorker
		if w == numWorkers-1 {
			endRow = p.Height
		}

		wg.Add(1)
		go func(startRow, endRow int) {
			defer wg.Done()
			for py := startRow; py < endRow; py++ {
				for px := 0; px < p.Width; px++ {
					z0 := julia.PixelToComplex(px, py, p.Width, p.Height, p)
					_, smooth := julia.Iterate(z0, p.C, p.MaxIter, p.EscapeRadius)
					buf[py*p.Width+px] = float32(smooth)
				}
			}
		}(startRow, endRow)
	}

	wg.Wait()
	return buf
}
