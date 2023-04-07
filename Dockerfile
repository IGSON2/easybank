# Build stage
FROM golang:1.20.0-alpine3.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o main .
RUN apk --no-cache add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.14.2
WORKDIR /app
# Copy the binary from the builder stage
COPY --from=builder /app/main . 
COPY --from=builder /app/migrate.linux-amd64 ./migrate 
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 4000
CMD ["./main"]

# wait-for를 통해 mysql:3360이 실행될 때까지 기다린다. 이후 /app/start.sh를 실행한다.
ENTRYPOINT [ "/app/wait-for.sh","mysql:3306","--","/app/start.sh" ]