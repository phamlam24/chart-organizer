# Backend Dockerfile
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Install required tools for code generation
RUN go install github.com/bufbuild/buf/cmd/buf@latest && \
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest

# Copy buf configuration files first
COPY buf.gen.yaml buf.gen.yaml
COPY buf.yaml buf.yaml
COPY contracts contracts

# Generate code from protobuf files (this will create backend/gen directory)
RUN buf generate

# Now set up the Go workspace in the backend subdirectory
WORKDIR /app/backend

# Copy go mod and sum files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy the backend source code
COPY backend/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Start a new stage from alpine
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

# Create storage directory
RUN mkdir -p storage

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/backend/main .

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