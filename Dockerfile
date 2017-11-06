FROM alpine

RUN apk add --no-cache ca-certificates

ADD matterbot /

ENTRYPOINT /matterbot
