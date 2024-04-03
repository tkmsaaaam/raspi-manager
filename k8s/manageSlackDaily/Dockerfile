# manage-slack-daily:latest
FROM alpine:latest as clone
RUN apk --update --no-cache add git
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/manage-slack.git --depth 1

FROM golang:1.21-alpine as build
COPY --from=clone /app/manage-slack /app/manage-slack
WORKDIR /app/manage-slack
RUN go mod tidy && go build daily.go

FROM alpine:latest
COPY --from=build /app/manage-slack/daily /app/daily
