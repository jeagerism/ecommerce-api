FROM golang:1.24-alpine

# Set working directory
WORKDIR /app

# Install git (for go mod dependencies)
RUN apk add --no-cache git

# Define build arguments
ARG DB_HOST
ARG DB_PORT
ARG DB_USER
ARG DB_PASSWORD
ARG DB_NAME
ARG SSL_MODE
ARG USER_SECRET
ARG SHOP_SECRET

# Make them available as environment variables at runtime
ENV DB_HOST=$DB_HOST
ENV DB_PORT=$DB_PORT
ENV DB_USER=$DB_USER
ENV DB_PASSWORD=$DB_PASSWORD
ENV DB_NAME=$DB_NAME
ENV SSL_MODE=$SSL_MODE
ENV USER_SECRET=$USER_SECRET
ENV SHOP_SECRET=$SHOP_SECRET

# Copy Go module files and download deps
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go binary
RUN go build -o ecommerce-api .

# Cloud Run requires app to listen on PORT (default: 8080)
ENV PORT=8080
EXPOSE 8080

# Run the binary
CMD ["./ecommerce-api"]
