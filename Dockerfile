FROM golang:1.24 AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy whole sources
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o supmap-gateway

# Build final image
FROM golang:1.24-alpine
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/supmap-gateway .

# Default values
ENV SUPMAP_GATEWAY_PORT=80
EXPOSE 80

ENTRYPOINT ["./supmap-gateway"]