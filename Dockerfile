FROM golang:1.18-alpine3.16 as build
WORKDIR /build
ADD . .
EXPOSE 8080
RUN go build -o rediswrapper .

FROM golang:1.18-alpine3.16 as test
WORKDIR /test
ADD . .
RUN apk add make git gcc musl-dev
CMD make _int-test

FROM alpine:3.16 as run
WORKDIR /app
COPY --from=build /build/rediswrapper .
CMD /app/rediswrapper
