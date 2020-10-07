FROM golang:latest as builder

WORKDIR /build
ADD . /build

RUN go get ./... && \
    go test ./... && \
    CGO_ENABLED=0 go build -o aws-s3-proxy ./cmd/

FROM alpine:latest
LABEL org.opencontainers.image.source=https://github.com/cirruslabs/aws-s3-proxy/

WORKDIR /svc
COPY --from=builder /build/aws-s3-proxy /svc/
ENTRYPOINT ["/svc/aws-s3-proxy"]
