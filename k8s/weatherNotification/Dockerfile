# weather-notification:latest
FROM debian:bookworm-slim AS clone
RUN apt-get update -y && apt-get install --no-install-recommends -y git ca-certificates
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/weather-notification.git --depth 1

FROM golang:1.24-bookworm AS build
COPY --from=clone /app/weather-notification /app/weather-notification
WORKDIR /app/weather-notification
RUN go mod tidy && go build

FROM debian:bookworm-slim
RUN apt-get update -y && apt-get install --no-install-recommends -y ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*
COPY --from=build /app/weather-notification/weather-notification /app/weather-notification
