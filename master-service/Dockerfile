FROM golang:alpine AS builder
WORKDIR /
ADD go.mod .
ADD go.sum .
RUN apk update && apk add --no-cache git tzdata ca-certificates gcc libtool make musl-dev protoc
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -o moyo-master-service
# FROM scratch
FROM alpine:latest
# COPY .env /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder moyo-master-service /
COPY --from=builder /database  /database
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# RUN ls -lah /
# RUN ls -lah /database/migrations/table
# Run the binary.
ENTRYPOINT ["/moyo-master-service"]
EXPOSE 8096