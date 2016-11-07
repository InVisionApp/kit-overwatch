FROM golang:1.6.3-alpine

# Install make
RUN apk update && apk add make git

ENV CGO_ENABLED=0\
    GOOS=linux

WORKDIR /go/src/github.com/InVisionApp/kit-overwatch
COPY . ./
