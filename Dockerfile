FROM golang:1.20.4-alpine3.18 AS base
COPY . /go/src/github.com/aschbacd/cloudflare-ddns
WORKDIR /go/src/github.com/aschbacd/cloudflare-ddns
RUN go build -a -tags netgo -ldflags '-w' -o /go/bin/cloudflare-ddns .

FROM alpine:3.18
RUN apk add tini bash
COPY --from=base /go/bin/cloudflare-ddns /usr/local/bin/cloudflare-ddns
WORKDIR /opt/cloudflare-ddns
ENTRYPOINT [ "tini", "bash", "-c" ]
CMD [ "while true; do cloudflare-ddns update; sleep 300; done" ]
