FROM golang:1.21.4 AS builder
WORKDIR /src
COPY . /src
RUN make build

FROM alpine:3.18.4
WORKDIR /findingway
COPY --from=builder /src/findingway .
COPY --from=builder /src/config.yaml .
ENTRYPOINT /findingway/findingway
