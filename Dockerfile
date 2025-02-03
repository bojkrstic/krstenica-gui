# Use the latest Go 1.22 image as the build stage
FROM golang:1.22-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create working directory inside container
WORKDIR /app

# Copy go.mod and go.sum first to download dependencies
COPY go.mod go.sum ./

# Ensure Go version is correct before running `go mod download`
RUN go version

# Download dependencies
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the application
# RUN go build -o app .
RUN CGO_ENABLED=0 GOOS=linux go build -o krstenica cmd/krstenica/main.go

# Use a minimal image for the final container
FROM alpine:latest 

WORKDIR /app
## We have to copy the output from our
## builder stage to our production stage
COPY --from=builder /app/krstenica .
COPY --from=builder /app/config . 


USER nobody

# Run the application
CMD ["./krstenica"]
