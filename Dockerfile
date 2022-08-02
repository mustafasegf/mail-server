# Build stage
FROM golang:1.18-alpine3.15 AS builder
WORKDIR /build
COPY . .
RUN go mod tidy
RUN go build -o main main.go

# Run stage
FROM alpine:3.15
WORKDIR /app
COPY --from=builder /build/main .
COPY .env .

CMD [ "/app/main" ]
