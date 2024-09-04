FROM golang:1.20 AS builder

WORKDIR /app

COPY main.go .

RUN go build -o arvanflux main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/arvanflux .

CMD ["./arvanflux"]
