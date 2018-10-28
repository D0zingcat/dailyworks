FROM alpine:latest

RUN apk add --update \
	bash \
	curl \
	openssl \
	&& rm -rf /var/cache/apk/*
COPY . .
WORKDIR dehydrated

ENTRYPOINT ["bash", "dehydrated"]

