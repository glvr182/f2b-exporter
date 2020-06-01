FROM golang:1.14 AS builder

WORKDIR /go/src/github.com/glvr182/f2b-exporter/
COPY . .
RUN CGO_ENABLED=0 go build -o /exe

FROM scratch
EXPOSE 8080
COPY --from=builder /exe /exe
ENTRYPOINT ["/exe"]