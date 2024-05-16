FROM golang:1.22.1 AS builder
WORKDIR /src
COPY . /src
RUN make build

FROM alpine:3.19.1
WORKDIR /findingway
COPY --from=builder /src/findingway .
COPY --from=builder /src/config.yaml .
ENTRYPOINT /findingway/findingway
