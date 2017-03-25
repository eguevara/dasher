FROM alpine:latest

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	ca-certificates

COPY . /go/src/github.com/eguevara/dasher

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		go \
		git \
		gcc \
		libc-dev \
		libgcc \
	&& cd /go/src/github.com/eguevara/dasher \
	&& go build -o /usr/bin/dasher . \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."
