FROM golang:1.11.1-alpine

RUN apk update && apk --no-cache add git

RUN addgroup -S clapper-bot && adduser -S -D clapper-bot clapper-bot

USER clapper-bot
WORKDIR "/home/clapper-bot"

# Install dep
RUN go get -u github.com/golang/dep/cmd/dep

COPY go.mod .
COPY Gopkg.toml .
COPY main.go .

# Ensure our dependencies are in the vendor directory.
RUN dep ensure

RUN go build .
CMD ["./clapper-bot"]
