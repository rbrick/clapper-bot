FROM golang:1.11.1-alpine

RUN apk update && apk --no-cache add git

RUN addgroup -S clapper-bot && adduser -S -D clapper-bot clapper-bot

USER clapper-bot
WORKDIR "/home/clapper-bot"

# Install dep
RUN go get -u github.com/golang/dep/cmd/dep

# Ensure our dependencies are in the vendor directory.
RUN dep ensure

COPY main.go .

RUN go build .
CMD ["./clapper-bot"]