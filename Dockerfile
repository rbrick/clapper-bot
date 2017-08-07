FROM golang:1.8.3-alpine

RUN apk update && apk --no-cache add git

RUN addgroup -S clapper-bot && adduser -S -D clapper-bot clapper-bot

USER clapper-bot
WORKDIR "/home/clapper-bot"

RUN go get gopkg.in/telegram-bot-api.v4

COPY main.go .

RUN go build .
CMD ["./clapper-bot"]