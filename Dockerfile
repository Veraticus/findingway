FROM golang:1.19.0 AS builder
WORKDIR /src
COPY . /src
RUN make build

FROM alpine:3.16.0
WORKDIR /trappingway
COPY --from=builder /src/trappingway .
COPY --from=builder /src/config.yaml .
ENTRYPOINT /trappingway/trappingway
