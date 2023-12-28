FROM golang:1.16-alpine3.12 as build
RUN apk update && apk add git
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/weather-notification.git
WORKDIR /app/weather-notification
RUN go get
RUN go build

FROM alpine:3.12
COPY --from=build /app/weather-notification/weather-notification /app/weather-notification
