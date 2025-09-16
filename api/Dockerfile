FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o app ./cmd/api

FROM alpine:latest AS run

WORKDIR /root/
COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./app"]
