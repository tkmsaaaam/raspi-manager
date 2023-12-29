# localhost:5000/manage-slack-daily
FROM alpine:3.12 as clone
RUN apk update && apk add git
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/manage-slack.git

FROM golang:1.21-alpine as build
COPY --from=clone /app/manage-slack /app/manage-slack
WORKDIR /app/manage-slack
RUN go mod tidy
RUN go build daily.go

FROM alpine:latest
COPY --from=build /app/manage-slack/daily /app/daily
