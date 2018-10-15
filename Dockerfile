FROM golang:1.11.1-alpine

RUN apk update && apk --no-cache add git curl

RUN addgroup -S clapper-bot && adduser -S -D clapper-bot clapper-bot

USER clapper-bot
WORKDIR "/home/clapper-bot"
RUN mkdir bin/

ENV GOPATH="/home/clapper-bot"
ENV GOBIN="${GOPATH}/bin"
ENV PATH="$PATH:${GOBIN}"

RUN mkdir ./src/
RUN mkdir ./src/clapper-bot
WORKDIR "./src/clapper-bot"

# Install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY go.mod .
COPY Gopkg.toml .
COPY main.go .

# Ensure our dependencies are in the vendor directory.
RUN dep ensure

RUN go build .
CMD ["./clapper-bot"]
