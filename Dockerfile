FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN /go/bin/swag init -g cmd/server/main.go

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -extldflags '-static'" \
    -o calculator \
    ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/calculator .
COPY --from=builder /app/docs ./docs

EXPOSE 8080
EXPOSE 50051

CMD ["./calculator"]
