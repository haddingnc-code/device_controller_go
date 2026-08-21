# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
ENTRYPOINT ["./main"]
# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
ENTRYPOINT ["./main"]
