# localhost:5000/weather-notification:latest
FROM alpine:latest as clone
RUN apk add --update --no-cache git
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/weather-notification.git --depth 1

FROM golang:1.21-alpine as build
COPY --from=clone /app/weather-notification /app/weather-notification
WORKDIR /app/weather-notification
RUN go mod tidy && go build

FROM alpine:latest
COPY --from=build /app/weather-notification/weather-notification /app/weather-notification
