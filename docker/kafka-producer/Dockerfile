FROM golang:1.12 as builder

COPY go-src/go.mod /src/
WORKDIR /src
RUN go mod download

COPY go-src /src
RUN go build -v -o kafka-producer.out ./cmd/kafka-producer/... || true

FROM ubuntu:18.04

COPY --from=builder /src/kafka-producer.out /usr/bin/.
CMD [ "kafka-producer" ]
