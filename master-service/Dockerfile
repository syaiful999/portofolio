FROM golang:alpine AS builder
WORKDIR /app
ADD go.mod go.sum ./
RUN apk update && apk add --no-cache git tzdata ca-certificates gcc libtool make musl-dev protoc
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -o moyo-master-service

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/moyo-master-service /moyo-master-service
COPY --from=builder /app/database/migrations /database/migrations
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENTRYPOINT ["/moyo-master-service"]
EXPOSE 8096
