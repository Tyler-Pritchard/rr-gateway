# Build Stage
FROM golang:1.23.1-alpine AS build

# Set working directory
WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the API Gateway binary
RUN go build -o api-gateway main.go

# Final Runtime Stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Install required dependencies
RUN apk add --no-cache ca-certificates curl

# Copy compiled binary from build stage
COPY --from=build /app/api-gateway .

# Expose API Gateway port
EXPOSE 8081

# Run the API Gateway
CMD ["./api-gateway"]
