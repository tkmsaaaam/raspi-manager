# rakuten-sec:latest
FROM debian:bookworm-slim as clone
RUN apt-get update -y && apt-get install --no-install-recommends -y git ca-certificates &&  apt-get clean && rm -rf /var/lib/apt/lists/*
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/browser-operator.git --depth 1

FROM golang:1.22-bookworm as build
ENV GOHOSTARCH=arm GOHOSTOS=linux GOARCH=arm GOOS=linux
RUN apt-get update -y  && apt-get install --no-install-recommends -y git ca-certificates &&  apt-get clean && rm -rf /var/lib/apt/lists/*
WORKDIR /work
RUN git clone https://github.com/tkmsaaaam/raspi-manager.git --depth 1
WORKDIR /work/raspi-manager/slackPublisher
RUN go install && go build

FROM node:22-bookworm-slim
COPY --from=clone /app/browser-operator/rakuten-sec /app/rakuten-sec
COPY --from=build /work/raspi-manager/slackPublisher/slackPublisher /app/slackPublisher
WORKDIR /app/rakuten-sec
RUN apt-get update -y  && apt-get install --no-install-recommends -y chromium ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*
ENV FILE_OUTPUT=true EXECUTABLE_PATH=/usr/bin/chromium
