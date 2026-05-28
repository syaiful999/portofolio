FROM --platform=linux/amd64 golang:1.20.14-alpine3.18

# Install protoc and dependencies
RUN apk add --no-cache protoc protobuf-dev

# Install Go protoc plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.16.2 && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.16.2

ENV PATH="/go/bin:${PATH}"
WORKDIR /source