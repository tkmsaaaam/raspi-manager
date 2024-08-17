# unitedcinemas:latest
FROM debian:bookworm-slim as build
WORKDIR /app
RUN sed -i 's@deb.debian.org/debian-security@ftp.riken.jp/Linux/debian/debian-security@g' /etc/apt/sources.list.d/debian.sources && sed -i 's@deb.debian.org/debian@ftp.jp.debian.org/debian@g' /etc/apt/sources.list.d/debian.sources
RUN apt-get update -y && apt-get install --no-install-recommends -y git ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*
RUN git clone https://github.com/tkmsaaaam/browser-operator.git --depth 1
WORKDIR /app/browser-operator/unitedcinemas

FROM golang:1.22-bookworm as build-publisher
ENV GOHOSTARCH=arm GOHOSTOS=linux GOARCH=arm GOOS=linux
RUN sed -i 's@deb.debian.org/debian-security@ftp.riken.jp/Linux/debian/debian-security@g' /etc/apt/sources.list.d/debian.sources && sed -i 's@deb.debian.org/debian@ftp.jp.debian.org/debian@g' /etc/apt/sources.list.d/debian.sources
RUN apt-get update -y  && apt-get install --no-install-recommends -y git ca-certificates &&  apt-get clean && rm -rf /var/lib/apt/lists/*
WORKDIR /work
RUN git clone https://github.com/tkmsaaaam/raspi-manager.git --depth 1
WORKDIR /work/raspi-manager/slackPublisher
RUN go install && go build

FROM node:22-bookworm-slim
COPY --from=build /app/browser-operator/unitedcinemas /app/unitedcinemas
COPY --from=build-publisher /work/raspi-manager/slackPublisher/slackPublisher /app/slackPublisher
RUN sed -i 's@deb.debian.org/debian-security@ftp.riken.jp/Linux/debian/debian-security@g' /etc/apt/sources.list.d/debian.sources && sed -i 's@deb.debian.org/debian@ftp.jp.debian.org/debian@g' /etc/apt/sources.list.d/debian.sources
RUN apt-get update -y  && apt-get install --no-install-recommends -y ca-certificates &&  apt-get clean && rm -rf /var/lib/apt/lists/*
WORKDIR /app/unitedcinemas
