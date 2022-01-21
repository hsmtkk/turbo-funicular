FROM golang:1.17 AS builder

WORKDIR /opt

COPY . .

RUN go build

FROM gcr.io/distroless/cc-debian11 AS runtime

COPY --from=builder /opt/turbo-funicular /usr/local/bin/turbo-funicular

EXPOSE 8000

ENTRYPOINT ["/usr/local/bin/turbo-funicular"]

