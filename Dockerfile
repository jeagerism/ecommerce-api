FROM golang:1.24-alpine

# Set working directory
WORKDIR /app

# Install git (for go mod dependencies)
RUN apk add --no-cache git

# Copy Go module files
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go binary
RUN go build -o ecommerce-api .

# Set the port Cloud Run expects
ENV PORT=8080
EXPOSE 8080

# Run the compiled binary
CMD ["./ecommerce-api"]
