# Use the official Golang image to build the application
FROM golang:1.23-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download all the dependencies.
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app

# Start a new stage from scratch to reduce image size
FROM alpine:latest  

# Install necessary dependencies (if required)
RUN apk --no-cache add ca-certificates curl

# Set environment variables
ENV GIN_MODE=release
ENV DB_HOST=host.docker.internal
ENV DB_USER=postgres
ENV DB_PASSWORD="powergrid@2?.net"
ENV DB_NAME=product_microservice
ENV DB_PORT=5432
ENV DB_SSLMODE=disable
ENV DB_TIMEZONE=UTC
ENV GRPC_PORT=50051

# Copy the compiled binary from the builder stage
COPY --from=builder /bin/app /bin/app

# Expose the port that your app will run on
EXPOSE 50051

# Add Healthcheck to verify if the app is healthy
HEALTHCHECK CMD curl --fail http://localhost:50051/health || exit 1

# Command to run the application
ENTRYPOINT ["/bin/app"]
