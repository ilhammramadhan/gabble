FROM golang:1.21-alpine AS builder

WORKDIR /app

# Cache bust - change this value to force rebuild
ARG CACHE_BUST=20260108_2103
RUN echo "Cache bust: $CACHE_BUST"

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
