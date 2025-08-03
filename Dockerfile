# Build stage
FROM golang:1.24-alpine AS builder

# Install necessary build tools
RUN apk add build-base

# Set environment variable
ENV ENVIRONMENT=production

# Set the current working directory inside the container
WORKDIR /build

# Copy the entire application code to the working directory
COPY . .

# Build the Go binary for the dashboard service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/bin ./cmd/

# Final stage
FROM alpine:latest

# Install CA certificates for HTTPS
RUN apk add ca-certificates

# Create a non-root user and group
RUN addgroup -g 3000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

# Create app directory and set ownership
RUN mkdir -p /app && chown -R appuser:appgroup /app

# Set the working directory to /app (not /root)
WORKDIR /app

# Copy the built binary from the builder stage to /app directory
COPY --from=builder /build/bin /app/bin

# Set execute permissions and ownership for the non-root user
RUN chmod +x /app/bin && chown appuser:appgroup /app/bin

# Expose the port your dashboard service listens on
EXPOSE 8080

# Switch to non-root user
USER appuser

# Run the binary (now located at /app/dashboard)
CMD ["/app/bin"]
