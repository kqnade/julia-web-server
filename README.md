# Julia Set Visualizer

A web application that renders Julia set fractals in real-time. Enter parameters in the browser, click Generate, and watch the fractal appear tile-by-tile.

## Quick Start

### Local development

```bash
mise install
go run .
# Open http://localhost:8080/satori/julia
```

### Docker

```bash
# ビルド＆起動（docker-compose）
docker compose up --build
# Open http://localhost:8080/satori/julia
```

ビルド済みイメージを直接使う場合：

```bash
# イメージのビルド
docker build --target prod -t julia-web-server:prod .

# コンテナの起動
docker run -p 8080:8080 julia-web-server:prod
# Open http://localhost:8080/satori/julia
```

デバッグ用イメージ（curl / telnet 付き）を使う場合：

```bash
docker build --target debug -t julia-web-server:debug .
docker run -p 8080:8080 julia-web-server:debug
```

## API

### `GET /satori/julia/api`

Computes a rectangular region of the Julia set and returns raw float32 binary data.

#### Required parameters

| Parameter | Format | Example | Description |
|---|---|---|---|
| `min_x` | float | `-2` | Real axis minimum |
| `max_x` | float | `2` | Real axis maximum |
| `min_y` | float | `-1.5` | Imaginary axis minimum |
| `max_y` | float | `1.5` | Imaginary axis maximum |
| `comp_const` | `real,imag` | `-0.7,0.27015` | Complex constant c |

#### Optional parameters

| Parameter | Range | Default | Description |
|---|---|---|---|
| `width` | 1-4096 | 256 | Output width in pixels |
| `height` | 1-4096 | 256 | Output height in pixels |
| `max_iter` | 1-10000 | 256 | Maximum iteration count |

#### Response

- **Success**: `Content-Type: application/octet-stream`
  - Body: `width * height` float32 values (little-endian)
  - `>= 0`: smooth iteration count (escaped point)
  - `-1.0`: interior point (did not escape)
- **Error**: `Content-Type: application/json`, Status 400
  - Body: `{"error": "reason"}`

#### Examples

```bash
# Normal: returns 262,144 bytes (256*256*4)
curl -o tile.bin "http://localhost:8080/satori/julia/api?min_x=-2&max_x=2&min_y=-1.5&max_y=1.5&comp_const=-0.7,0.27015"

# Error: missing parameter
curl -s "http://localhost:8080/satori/julia/api?min_x=-2"
# {"error":"missing required parameter: max_x"}
```

## Algorithm

### Julia Set Iteration

For each pixel mapped to complex coordinate `z0`:

```
z = z0
for i = 0..max_iter:
    if |z|² > escape_radius²:
        return escaped with smooth value
    z = z² + c
return not escaped (-1.0)
```

- **Escape radius**: 2.0 (mathematically proven: if |z| > 2, the sequence diverges)
- **Smooth coloring**: `i + 1 - log(log(|z|)) / log(2)` — logarithmic interpolation eliminates banding artifacts
- **Optimization**: Compare `|z|²` instead of `|z|` to avoid sqrt per iteration

### HSV Coloring (frontend)

- Escaped points: `Hue = (smooth * 10) mod 360`, Saturation = 1.0, Value = 1.0
- Interior points: black (0, 0, 0)

### Tile-based Rendering

The 800x600 canvas is divided into 256x256 tiles, fetched in parallel for progressive display.

## Tests

```bash
go test ./...           # Run all tests
go test -race ./...     # With race condition detection
go test -cover ./...    # With coverage
go test -v ./...        # Verbose output
```

## Project Structure

```
├── main.go                     # Server entry point, routing, embed
├── internal/
│   ├── julia/julia.go          # Core iteration math
│   ├── renderer/renderer.go    # Parallel float32 buffer generation
│   └── handler/
│       ├── handler.go          # HTTP handler
│       └── params.go           # Query parameter parsing/validation
├── web/
│   ├── index.html              # UI (form + canvas)
│   └── app.js                  # Tile splitting, fetch, HSV coloring
├── Dockerfile                  # Multi-stage build (builder/debug/prod)
└── docker-compose.yml          # One-command startup
```
