################################
# STEP 1 build executable binary
################################

# https://songrgg.github.io/operation/how-to-build-a-smallest-docker-image/

FROM golang:alpine3.12 as builder

ENV GOPATH=/go

RUN apk update && apk add alpine-sdk git bash && rm -rf /var/cache/apk/*

WORKDIR /go/src/github.com/bzimmer/gravl

COPY go.mod ./
COPY go.sum ./

RUN go mod download

ADD pkg pkg
ADD cmd cmd

# RUN go test ./...

ARG VERSION
RUN go build -o dist/gravl -ldflags "-X github.com/bzimmer/gravl/pkg.BuildVersion=$VERSION -X github.com/bzimmer/gravl/pkg.BuildTime=`date +%Y-%d-%mT%H:%M:%S`" cmd/gravl/*.go

##############################
# STEP 2 build a smaller image
##############################

FROM alpine:3.12.0

# https://github.com/googleapis/google-cloud-go/issues/928
RUN apk --no-cache --update add ca-certificates

WORKDIR /app

COPY --from=builder /go/src/github.com/bzimmer/gravl/dist/gravl .

ENV GRAVL_PORT=8080
ENV GIN_MODE=debug

ENTRYPOINT ["/app/gravl", "serve"]
