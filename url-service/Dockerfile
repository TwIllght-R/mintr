# Use Golang base image for building
FROM golang:1.24 AS builder

# Set Go modules to be on by default
ENV GO111MODULE=on

# Set the working directory to /app
WORKDIR /app

RUN go install github.com/air-verse/air@latest


ENV PATH="/go/bin:$PATH"

ADD ./ .

# Download dependencies to cache them
RUN go mod download


# Build the Go binary
RUN go build -o /go/bin/app.bin
