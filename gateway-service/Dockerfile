FROM golang:alpine AS builder
WORKDIR /
ADD go.mod .
ADD go.sum .
RUN apk update && apk add --no-cache git tzdata ca-certificates gcc libtool make musl-dev protoc
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -o moyo-gateway-service
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.11 && \
    wget -qO /bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe
# FROM scratch
FROM alpine
# COPY .env /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder moyo-gateway-service /
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /bin/grpc_health_probe /usr/local/bin/grpc_health_probe
# Run the binary.
ENTRYPOINT ["/moyo-gateway-service"]
EXPOSE 8080