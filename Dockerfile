FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o ecommerce-api .

EXPOSE 8080

CMD ["./ecommerce-api"]
