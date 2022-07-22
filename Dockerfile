ARG ALPINE_VERSION=3.16
ARG GO_VERSION=1.18

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as build
WORKDIR /build
ADD . .
EXPOSE 8080
RUN go build -o rediswrapper .

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as test
WORKDIR /test
ADD . .
RUN apk add make git gcc musl-dev
CMD make _int-test

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as fuzz
WORKDIR /fuzz
ADD . .
RUN apk add make git gcc musl-dev
CMD go test -v -fuzz FuzzIntegration

FROM alpine:${ALPINE_VERSION} as run
WORKDIR /app
COPY --from=build /build/rediswrapper .
CMD /app/rediswrapper
