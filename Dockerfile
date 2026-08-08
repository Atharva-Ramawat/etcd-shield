FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /etcd-shield ./cmd/etcd-shield

FROM alpine:latest
WORKDIR /app
COPY --from=builder /etcd-shield /etcd-shield
COPY --from=builder /app/config.yaml ./config.yaml
ENTRYPOINT ["/etcd-shield"]