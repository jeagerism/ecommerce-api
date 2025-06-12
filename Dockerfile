# Start from a lightweight Go base image
FROM golang:1.22-alpine

# Create working directory inside container
WORKDIR /app

# Install git and other dependencies (for go modules)
RUN apk add --no-cache git

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy rest of the source code
COPY . .

# Build the Go app
RUN go build -o ecommerce-api .

# Run the binary
CMD ["./ecommerce-api"]
