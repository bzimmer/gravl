################################
# STEP 1 build executable binary
################################

# https://songrgg.github.io/operation/how-to-build-a-smallest-docker-image/

FROM golang:alpine3.12 as builder

ENV GOPATH=/go

RUN apk update && apk add alpine-sdk git bash && rm -rf /var/cache/apk/*

WORKDIR /go/src/github.com/bzimmer/wta

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY main.go ./
ADD pkg pkg
ADD cmd cmd

ADD testdata testdata
RUN go test -v ./...

ARG VERSION
RUN go build -o dist/atk -ldflags "-X github.com/bzimmer/wta/pkg.BuildVersion=$VERSION" main.go

##############################
# STEP 2 build a smaller image
##############################

FROM alpine:3.12.0

# https://github.com/googleapis/google-cloud-go/issues/928
RUN apk --no-cache --update add ca-certificates

WORKDIR /app

COPY --from=builder /go/src/github.com/bzimmer/wta/dist/atk .

ENV ATK_PORT=8080
ENV GIN_MODE=release

ENTRYPOINT ["/app/atk", "serve"]
