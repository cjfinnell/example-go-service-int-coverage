FROM golang:1.18-alpine3.16 as build
WORKDIR /build
ADD . .
EXPOSE 8080
RUN go build -o rediswrapper .

FROM alpine:3.16 as run
WORKDIR /app
COPY --from=build /build/rediswrapper .
CMD /app/rediswrapper
