# ==========================================
# 1. Builder: Build the Go binary
# ==========================================
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o julia-server .

# ==========================================
# 2. Debug: Development/debugging image
# ==========================================
FROM alpine:latest AS debug
WORKDIR /app
RUN apk add --no-cache curl busybox-extras
COPY --from=builder /app/julia-server .
EXPOSE 8080
CMD ["./julia-server"]

# ==========================================
# 3. Prod: Minimal production image
# ==========================================
FROM scratch AS prod
COPY --from=builder /app/julia-server /julia-server
EXPOSE 8080
ENTRYPOINT ["/julia-server"]
