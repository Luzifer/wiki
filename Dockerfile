FROM golang:alpine as builder

COPY . /go/src/github.com/Luzifer/wiki
WORKDIR /go/src/github.com/Luzifer/wiki

RUN set -ex \
 && apk --no-cache add \
      curl \
      git \
      make \
      nodejs \
      npm \
 && make frontend \
 && go build \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly \
      -o /go/bin/wiki

RUN set -ex \
 && curl -sSfLo /usr/local/bin/dumb-init "https://github.com/Yelp/dumb-init/releases/download/v1.2.2/dumb-init_1.2.2_amd64" \
 && curl -sSfLo /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/1.11/gosu-amd64" \
 && chmod +x \
      /usr/local/bin/dumb-init \
      /usr/local/bin/gosu


FROM alpine:latest

ENV DATA_DIR=/data

LABEL maintainer "Knut Ahlers <knut@ahlers.me>"

RUN set -ex \
 && apk --no-cache add \
      bash \
      ca-certificates \
 && adduser -D -h /home/wiki -S -u 1000 wiki

COPY                docker-entrypoint.sh      /usr/local/bin/docker-entrypoint
COPY --from=builder /go/bin/wiki              /usr/local/bin/wiki
COPY --from=builder /usr/local/bin/dumb-init  /usr/local/bin/dumb-init
COPY --from=builder /usr/local/bin/gosu       /usr/local/bin/gosu

EXPOSE 3000
VOLUME ["/data"]

ENTRYPOINT ["/usr/local/bin/docker-entrypoint"]
CMD ["--"]

# vim: set ft=Dockerfile:
