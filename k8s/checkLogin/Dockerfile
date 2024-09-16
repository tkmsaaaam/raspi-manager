# check-login:latest
FROM golang:1.22-bookworm AS build
ENV GOHOSTARCH=arm GOHOSTOS=linux GOARCH=arm GOOS=linux
RUN apt-get update -y  && apt-get install --no-install-recommends -y git &&  apt-get clean && rm -rf /var/lib/apt/lists/*
WORKDIR /work
RUN git clone https://github.com/tkmsaaaam/raspi-manager.git --depth 1
WORKDIR /work/raspi-manager/slackPublisher
RUN go install && go build

FROM debian:bookworm-slim
COPY AuthlogParser/bin/Release/net8.0/linux-arm/publish/ /app/
COPY --from=build /work/raspi-manager/slackPublisher/slackPublisher /app/slackPublisher
RUN apt-get update -y  && apt-get install --no-install-recommends -y ca-certificates tzdata jq libc6 libgcc-s1 libgssapi-krb5-2 libicu72 libstdc++6 libssl3 zlib1g && apt-get clean && rm -rf /var/lib/apt/lists/*
WORKDIR /app
CMD ["/bin/sh", "-c", "/app/AuthlogParser"]
