# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Runtime Stage
FROM alpine:latest

WORKDIR /app

# Copy binary
COPY --from=builder /app/main .

# Copy migrations
COPY --from=builder /app/scripts/migrations ./scripts/migrations

EXPOSE 8080

CMD ["./main"]