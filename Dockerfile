FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/app/main.go

# server
FROM alpine:latest

RUN adduser -D appuser
USER appuser

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8001

CMD ["./server"]

# docker build -t jod .
# docker run jod:latest