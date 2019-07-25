FROM golang:1.12 as builder

COPY go-src/go.mod /src/
WORKDIR /src
RUN go mod download

COPY go-src /src
RUN go build -v -o kafka-consumergroup.out ./cmd/kafka-consumergroup/... || true

FROM ubuntu:18.04

COPY --from=builder /src/kafka-consumergroup.out /usr/bin/.
CMD [ "kafka-consumergroup" ]
