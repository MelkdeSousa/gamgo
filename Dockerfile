# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:3.17

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy any config files if needed
COPY rawg-open-api.json .

EXPOSE 3000

CMD ["./main"]