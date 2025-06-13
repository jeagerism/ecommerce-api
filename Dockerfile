# Use Go official image with Alpine for small size
FROM golang:1.21-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum and download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go application
RUN go build -o ecommerce-api .

# Final stage: use minimal base image
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Install necessary SSL certs and copy binary
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/ecommerce-api .

# Expose the port Cloud Run expects
EXPOSE 8080

# Run the Go application
CMD ["./ecommerce-api"]
