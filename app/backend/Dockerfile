# ==================================
# STAGE 1: builder (Tahap Kompilasi)
# ==================================
FROM golang:1.25-alpine AS builder

# 1. Tentukan working directory
WORKDIR /app

# 2. Salin go.mod dan go.sum
# Karena context build = . (backend/), kita bisa langsung copy ke /app
COPY go.mod .
COPY go.sum .

# 3. Unduh dependencies
RUN go mod download

# 4. Salin seluruh kode sumber Golang
COPY . .

# 5. Kompilasi aplikasi
# main.go ada di root folder /app
RUN go build -o go-app .


# ==================================
# STAGE 2: final (Tahap Produksi)
# ==================================
FROM alpine:latest

WORKDIR /root/

# Copy executable yang sudah dikompilasi
COPY --from=builder /app/go-app .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./go-app"]