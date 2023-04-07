# Build stage
FROM golang:1.20.0-alpine3.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

# Run stage
FROM alpine:3.14.2
WORKDIR /app
# Copy the binary from the builder stage
COPY --from=builder /app/main . 
COPY app.env .
EXPOSE 4000
CMD ["./main"]