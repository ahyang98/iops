FROM golang:1.21.5 AS builder
WORKDIR /src
COPY tcp_server .
RUN CGO_ENABLED=0 go build -o /bin/tcp_server ./tcp_server.go

FROM alpine:3.20.2
COPY --from=builder /bin/tcp_server /bin/tcp_server
    
CMD ["/bin/tcp_server"]