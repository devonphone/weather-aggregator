# Step 1: Build the Go app
FROM golang:1.20-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download the Go dependencies
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Copy the .env file into the container
COPY .env .env

# Build the Go app
RUN go build -o main .

# Step 2: Run the Go app
FROM alpine:latest

# Install ca-certificates (needed for Redis over HTTPS) and bash (optional)
RUN apk --no-cache add ca-certificates bash

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy the .env file
COPY --from=builder /app/.env .env

# Expose the port the app runs on
EXPOSE 8080

# Run the Go app
CMD ["./main"]