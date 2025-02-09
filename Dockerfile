FROM golang:1.23.6-alpine AS builder

WORKDIR /app

# Copy go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build the application
COPY . .
RUN go build -o main ./cmd/app

#Lightweight container
FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .

# Expose the app port at 8080
EXPOSE 8080

# Run the binary
CMD ["./main"]
