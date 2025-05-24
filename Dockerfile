# Use an official Golang runtime as a parent image
FROM golang:1.24-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Build the Go app
RUN go build -o server ./cmd/server
