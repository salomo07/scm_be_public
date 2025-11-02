# Stage 1: Build Golang Application
FROM golang:alpine AS builder
WORKDIR /app

# Copy go.mod dan go.sum sebelum kode lainnya untuk memanfaatkan layer cache
# RUN go mod download
COPY go.mod go.sum ./

# Copy semua kode setelah dependensi terunduh
COPY . .

# Build aplikasi
RUN go mod tidy && go build -o main

# Stage 2: Final Container
FROM alpine:latest
WORKDIR /root/

# Install dependensi yang diperlukan dalam runtime
RUN apk add --no-cache ca-certificates

# Copy binary yang sudah di-build
COPY --from=builder /app/main /root/main

# Expose port aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["/root/main"]
