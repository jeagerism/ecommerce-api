# --- Stage 1: Builder Stage ---
    FROM golang:1.22-alpine AS builder

    # Set working directory inside the container
    WORKDIR /app
    
    # Install git (required by go mod download for private modules, if any)
    # Note: This is only in the builder stage, not the final image.
    RUN apk add --no-cache git
    
    # Copy Go module files and download dependencies
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy source code
    COPY . .
    
    # Build the Go binary
    # CGO_ENABLED=0 creates a statically-linked binary, removing glibc dependency for Alpine
    # -o specifies the output binary name
    RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s -w' -o ecommerce-api .
    
    # --- Stage 2: Final (Runtime) Stage ---
    # Use a minimal base image to run the compiled binary
    FROM alpine:latest
    
    # Set working directory inside the container
    WORKDIR /app
    
    # Copy the compiled binary from the builder stage
    COPY --from=builder /app/ecommerce-api .
    
    # Cloud Run requires the application to listen on the PORT environment variable
    # The Go application already reads this from os.Getenv("PORT")
    # It's good practice to EXPOSE the port, even though Cloud Run sets it.
    EXPOSE 8080
    
    # Command to run the application when the container starts
    CMD ["./ecommerce-api"]
    
    # No ARG or ENV for DB_HOST, DB_PASSWORD, etc. in Dockerfile!
    # These will be passed at runtime by Cloud Run.