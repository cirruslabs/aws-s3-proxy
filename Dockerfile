FROM golang:1.11-alpine as builder

ENV CGO_ENABLED 0

RUN apk update && apk upgrade && \
    apk add --no-cache git

WORKDIR /tmp/aws-s3-proxy
ADD . /tmp/aws-s3-proxy

RUN go get -t -v ./... && \
    go test -v ./... && \
    go build -o aws-s3-proxy ./cmd/

FROM alpine

RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates

COPY --from=builder /tmp/aws-s3-proxy/aws-s3-proxy /bin/
ENTRYPOINT ["aws-s3-proxy"]