package handler

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/kqnade/julia-web-server/internal/julia"
)

const (
	defaultWidth   = 256
	defaultHeight  = 256
	defaultMaxIter = julia.DefaultMaxIter

	minDimension = 1
	maxDimension = 4096
	minMaxIter   = 1
	maxMaxIter   = 10000
)

// parseParams parses and validates query parameters, returning julia.Params or an error message.
func parseParams(q url.Values) (julia.Params, string) {
	// Required parameters
	minXStr := q.Get("min_x")
	if minXStr == "" {
		return julia.Params{}, "missing required parameter: min_x"
	}
	maxXStr := q.Get("max_x")
	if maxXStr == "" {
		return julia.Params{}, "missing required parameter: max_x"
	}
	minYStr := q.Get("min_y")
	if minYStr == "" {
		return julia.Params{}, "missing required parameter: min_y"
	}
	maxYStr := q.Get("max_y")
	if maxYStr == "" {
		return julia.Params{}, "missing required parameter: max_y"
	}
	compConstStr := q.Get("comp_const")
	if compConstStr == "" {
		return julia.Params{}, "missing required parameter: comp_const"
	}

	// Parse required float parameters
	minX, err := strconv.ParseFloat(minXStr, 64)
	if err != nil || math.IsNaN(minX) || math.IsInf(minX, 0) {
		return julia.Params{}, fmt.Sprintf("invalid min_x: %q is not a valid number", minXStr)
	}
	maxX, err := strconv.ParseFloat(maxXStr, 64)
	if err != nil || math.IsNaN(maxX) || math.IsInf(maxX, 0) {
		return julia.Params{}, fmt.Sprintf("invalid max_x: %q is not a valid number", maxXStr)
	}
	minY, err := strconv.ParseFloat(minYStr, 64)
	if err != nil || math.IsNaN(minY) || math.IsInf(minY, 0) {
		return julia.Params{}, fmt.Sprintf("invalid min_y: %q is not a valid number", minYStr)
	}
	maxY, err := strconv.ParseFloat(maxYStr, 64)
	if err != nil || math.IsNaN(maxY) || math.IsInf(maxY, 0) {
		return julia.Params{}, fmt.Sprintf("invalid max_y: %q is not a valid number", maxYStr)
	}

	// Parse comp_const
	parts := strings.SplitN(compConstStr, ",", 3)
	if len(parts) != 2 {
		return julia.Params{}, fmt.Sprintf("invalid comp_const: %q must be two comma-separated numbers", compConstStr)
	}
	cReal, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil || math.IsNaN(cReal) || math.IsInf(cReal, 0) {
		return julia.Params{}, fmt.Sprintf("invalid comp_const real part: %q is not a valid number", parts[0])
	}
	cImag, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil || math.IsNaN(cImag) || math.IsInf(cImag, 0) {
		return julia.Params{}, fmt.Sprintf("invalid comp_const imaginary part: %q is not a valid number", parts[1])
	}

	// Validate ranges
	if minX >= maxX {
		return julia.Params{}, fmt.Sprintf("min_x (%v) must be less than max_x (%v)", minX, maxX)
	}
	if minY >= maxY {
		return julia.Params{}, fmt.Sprintf("min_y (%v) must be less than max_y (%v)", minY, maxY)
	}

	// Optional parameters with defaults
	width := defaultWidth
	if ws := q.Get("width"); ws != "" {
		w, err := strconv.Atoi(ws)
		if err != nil {
			return julia.Params{}, fmt.Sprintf("invalid width: %q is not a valid integer", ws)
		}
		if w < minDimension || w > maxDimension {
			return julia.Params{}, fmt.Sprintf("width must be between %d and %d, got %d", minDimension, maxDimension, w)
		}
		width = w
	}

	height := defaultHeight
	if hs := q.Get("height"); hs != "" {
		h, err := strconv.Atoi(hs)
		if err != nil {
			return julia.Params{}, fmt.Sprintf("invalid height: %q is not a valid integer", hs)
		}
		if h < minDimension || h > maxDimension {
			return julia.Params{}, fmt.Sprintf("height must be between %d and %d, got %d", minDimension, maxDimension, h)
		}
		height = h
	}

	maxIter := defaultMaxIter
	if ms := q.Get("max_iter"); ms != "" {
		m, err := strconv.Atoi(ms)
		if err != nil {
			return julia.Params{}, fmt.Sprintf("invalid max_iter: %q is not a valid integer", ms)
		}
		if m < minMaxIter || m > maxMaxIter {
			return julia.Params{}, fmt.Sprintf("max_iter must be between %d and %d, got %d", minMaxIter, maxMaxIter, m)
		}
		maxIter = m
	}

	return julia.Params{
		MinX:         minX,
		MaxX:         maxX,
		MinY:         minY,
		MaxY:         maxY,
		C:            complex(cReal, cImag),
		Width:        width,
		Height:       height,
		MaxIter:      maxIter,
		EscapeRadius: julia.DefaultEscapeRadius,
	}, ""
}
