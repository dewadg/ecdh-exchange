FROM golang:1.20-alpine AS builder

WORKDIR /build

COPY . .

RUN go mod download
RUN go build -o ecdh-exchange cmd/server/cmd_server.go

FROM alpine

WORKDIR /usr/bin

COPY --from=builder /build/ecdh-exchange .
RUN chmod a+x ecdh-exchange

ENV APP_ENV production

CMD ["/usr/bin/ecdh-exchange"]
