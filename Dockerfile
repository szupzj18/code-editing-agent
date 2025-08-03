# Multi-stage build for Go application
FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application (ignore vendor directory to avoid vendoring inconsistency)
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/main .

# Expose port (if needed)
EXPOSE 8080

# Run application
CMD ["./main"]
