FROM golang:1.16-alpine3.12 as build
RUN apk update && apk add git
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/manage-slack.git
WORKDIR /app/manage-slack
RUN go get
RUN go build main.go

FROM alpine:3.12
COPY --from=build /app/manage-slack/main /app/main
