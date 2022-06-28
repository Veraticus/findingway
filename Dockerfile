FROM golang:1.18.1 AS builder
WORKDIR /src
COPY . /src
RUN make build

FROM alpine:3.16.0
WORKDIR /trappingway
COPY --from=builder /src/trappingway .
ENTRYPOINT /trappingway/trappingway
