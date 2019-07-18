FROM golang:1.12 as builder

COPY go.mod /src/
WORKDIR /src
RUN go mod download

COPY . /src
RUN go build -v ./cmd/hello-world/...

FROM ubuntu:18.04

COPY --from=builder /src/hello-world /usr/bin/.
CMD [ "hello-world" ]