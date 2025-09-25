# Backend Dockerfile
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Install required tools for code generation
RUN go install github.com/bufbuild/buf/cmd/buf@latest && \
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest

# Copy go mod and sum files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code (including contracts and buf config files from root)
COPY backend/ .
COPY contracts contracts
COPY buf.gen.yaml buf.gen.yaml
COPY buf.yaml buf.yaml

# Generate code from protobuf files
RUN buf generate

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Start a new stage from alpine
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

# Create storage directory
RUN mkdir -p storage

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Create a basic .env file if it doesn't exist
RUN echo "JWT_KEY=your-super-secret-jwt-key-change-me-in-production" > .env

# Set environment variables
ENV PORT=8080
ENV DB_PATH=./storage/chart-organizer.db
ENV DATASET_STORAGE_PATH=./storage/datasets

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]