FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /myapp ./cmd/server

FROM alpine:3.23

WORKDIR /

RUN apk --no-cache add ca-certificates

COPY --from=builder /myapp /myapp

EXPOSE 8080

ENTRYPOINT ["/myapp"]
