# Use the official Golang base image with Alpine Linux
FROM golang:1.16-alpine

# Set the working directory inside the container
WORKDIR /app

ENV CONF_DIR=/app/conf

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the Go code
COPY . .

# Build the Go application
RUN go build -o main .

# Set the container entrypoint
ENTRYPOINT [ "./main" ]
