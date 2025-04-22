# Build stage
FROM golang:1.24 AS builder

WORKDIR /app    
COPY . .

RUN go mod tidy

# Build statically-linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64    go build -o app.bin .

# Final minimal stage
FROM alpine:3.21

WORKDIR /root
COPY --from=builder /app/app.bin .

EXPOSE 8080

CMD ["./app.bin"]
