#
# builder
#
FROM alpine as builder


# musl-dev is necessary to compile a
# golang executable
RUN apk add --no-cache go git musl-dev

WORKDIR /go/src/github.com/section77/matterbot

COPY . .

ENV GOPATH=/go

RUN go get

ENV CGO_ENABLE=0

RUN go install -ldflags "-X main.version=$(git describe --tags --always)"


#
# runtime
#

FROM alpine

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/matterbot /

ENTRYPOINT /matterbot
