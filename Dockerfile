FROM golang:1.15.3-alpine3.12 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/neoway-case

COPY go.mod go.sum ./
COPY util util
COPY db db
COPY schema schema
COPY configuration configuration
COPY errors errors
COPY consumption-service consumption-service

RUN GO111MODULE=on go install ./consumption-service

FROM alpine:3.12
WORKDIR /usr/bin
COPY --from=build /go/bin .