FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git
WORKDIR /src

COPY go.mod go.sum ./
COPY ./internal ./internal
COPY ./cmd ./cmd

RUN go build -o /app/kubewatch-backend ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/kubewatch-backend /usr/local/bin/kubewatch-backend
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/kubewatch-backend"]
