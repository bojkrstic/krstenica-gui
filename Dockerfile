# 1. Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Go mod fajlovi
COPY go.mod go.sum ./
RUN go mod download

# Ceo kod (uključuje web/templates, internal, cmd, itd.)
COPY . .

# Build binara (main je u cmd/krstenica)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/krstenica

# 2. Runtime stage
FROM alpine:latest

WORKDIR /app

# Binarna aplikacija
COPY --from=builder /app/server .

# 🔹 OVO JE KLJUČNO: kopiramo web/ da bi /app/web/templates postojao
COPY web /app/web

# 🔹 dodaj slike koje PDF i štampa koriste
COPY pictures /app/pictures

EXPOSE 8011

CMD ["./server"]
