package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const validQuery = "min_x=-2&max_x=2&min_y=-1.5&max_y=1.5&comp_const=-0.7,0.27015"

func TestJuliaAPI_Success(t *testing.T) {
	req := httptest.NewRequest("GET", "/satori/julia/api?"+validQuery, nil)
	w := httptest.NewRecorder()

	JuliaAPI(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/octet-stream" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/octet-stream")
	}
	// Default 256×256, 4 bytes per float32
	wantSize := 256 * 256 * 4
	if w.Body.Len() != wantSize {
		t.Errorf("body size = %d, want %d", w.Body.Len(), wantSize)
	}
}

func TestJuliaAPI_CustomSize(t *testing.T) {
	req := httptest.NewRequest("GET", "/satori/julia/api?"+validQuery+"&width=64&height=32", nil)
	w := httptest.NewRecorder()

	JuliaAPI(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	wantSize := 64 * 32 * 4
	if w.Body.Len() != wantSize {
		t.Errorf("body size = %d, want %d", w.Body.Len(), wantSize)
	}
}

func TestJuliaAPI_ValidationErrors(t *testing.T) {
	tests := []struct {
		name            string
		query           string
		wantErrContains string
	}{
		{"missing min_x", "max_x=2&min_y=-1.5&max_y=1.5&comp_const=-0.7,0.27015", "min_x"},
		{"missing max_x", "min_x=-2&min_y=-1.5&max_y=1.5&comp_const=-0.7,0.27015", "max_x"},
		{"missing min_y", "min_x=-2&max_x=2&max_y=1.5&comp_const=-0.7,0.27015", "min_y"},
		{"missing max_y", "min_x=-2&max_x=2&min_y=-1.5&comp_const=-0.7,0.27015", "max_y"},
		{"missing comp_const", "min_x=-2&max_x=2&min_y=-1.5&max_y=1.5", "comp_const"},
		{"invalid comp_const format", "min_x=-2&max_x=2&min_y=-1.5&max_y=1.5&comp_const=abc", "comp_const"},
		{"comp_const one part invalid", "min_x=-2&max_x=2&min_y=-1.5&max_y=1.5&comp_const=1.0,abc", "comp_const"},
		{"min_x > max_x", "min_x=2&max_x=-2&min_y=-1.5&max_y=1.5&comp_const=-0.7,0.27015", "min_x"},
		{"min_y > max_y", "min_x=-2&max_x=2&min_y=1.5&max_y=-1.5&comp_const=-0.7,0.27015", "min_y"},
		{"min_x == max_x", "min_x=0&max_x=0&min_y=-1.5&max_y=1.5&comp_const=-0.7,0.27015", "min_x"},
		{"max_iter too low", validQuery + "&max_iter=0", "max_iter"},
		{"max_iter too high", validQuery + "&max_iter=99999", "max_iter"},
		{"width too low", validQuery + "&width=0", "width"},
		{"width too high", validQuery + "&width=5000", "width"},
		{"height too low", validQuery + "&height=0", "height"},
		{"height too high", validQuery + "&height=5000", "height"},
		{"min_x not a number", "min_x=abc&max_x=2&min_y=-1.5&max_y=1.5&comp_const=-0.7,0.27015", "min_x"},
		{"min_x is NaN", "min_x=NaN&max_x=2&min_y=-1.5&max_y=1.5&comp_const=-0.7,0.27015", "min_x"},
		{"max_x is Inf", "min_x=-2&max_x=Inf&min_y=-1.5&max_y=1.5&comp_const=-0.7,0.27015", "max_x"},
		{"min_y is -Inf", "min_x=-2&max_x=2&min_y=-Inf&max_y=1.5&comp_const=-0.7,0.27015", "min_y"},
		{"max_y is NaN", "min_x=-2&max_x=2&min_y=-1.5&max_y=NaN&comp_const=-0.7,0.27015", "max_y"},
		{"comp_const real is NaN", "min_x=-2&max_x=2&min_y=-1.5&max_y=1.5&comp_const=NaN,0.27015", "comp_const"},
		{"comp_const imag is Inf", "min_x=-2&max_x=2&min_y=-1.5&max_y=1.5&comp_const=-0.7,+Inf", "comp_const"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/satori/julia/api?"+tt.query, nil)
			w := httptest.NewRecorder()

			JuliaAPI(w, req)

			resp := w.Result()
			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
			}
			if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want %q", ct, "application/json")
			}

			var body map[string]string
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}
			if !strings.Contains(body["error"], tt.wantErrContains) {
				t.Errorf("error = %q, want containing %q", body["error"], tt.wantErrContains)
			}
		})
	}
}
