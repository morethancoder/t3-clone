# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build both applications
RUN CGO_ENABLED=0 GOOS=linux go build -o t3-clone ./cmd/t3-clone/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o db-server ./cmd/db-server/main.go

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy built binaries
COPY --from=builder /app/t3-clone .
COPY --from=builder /app/db-server .

# Copy static files and other necessary directories
COPY --from=builder /app/static ./static
COPY --from=builder /app/views ./views
COPY --from=builder /app/tailwindcss ./tailwindcss
COPY --from=builder /app/pb_schema.json ./pb_schema.json

# Create pb_data directory
RUN mkdir -p pb_data

# Set environment variables with container-appropriate defaults
ENV ENV=production
ENV DB_URL=http://localhost:8090
ENV APP_URL=http://localhost:8080
ENV OPENROUTER_API_KEY=your_api_key_here

# Expose ports
EXPOSE 8080 8090

# Create a startup script
RUN echo '#!/bin/sh' > start.sh &&     echo './db-server serve &' >> start.sh &&     echo 'sleep 2' >> start.sh &&     echo './t3-clone' >> start.sh &&     chmod +x start.sh

CMD ["./start.sh"]
