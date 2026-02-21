package handler

import (
	"encoding/binary"
	"encoding/json"
	"net/http"

	"github.com/kqnade/julia-web-server/internal/renderer"
)

// JuliaAPI handles GET requests to compute Julia set tiles.
func JuliaAPI(w http.ResponseWriter, r *http.Request) {
	params, errMsg := parseParams(r.URL.Query())
	if errMsg != "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": errMsg})
		return
	}

	buf := renderer.Render(params)

	w.Header().Set("Content-Type", "application/octet-stream")
	binary.Write(w, binary.LittleEndian, buf)
}
