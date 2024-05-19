# rakuten-sec:latest
FROM debian:bookworm-slim as clone
RUN apt-get update -y && apt-get install --no-install-recommends -y git ca-certificates
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/browser-operator.git --depth 1

FROM node:21-bookworm-slim
COPY --from=clone /app/browser-operator/rakuten-sec /app/rakuten-sec
WORKDIR /app/rakuten-sec
RUN apt-get update -y  && apt-get install --no-install-recommends -y curl ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*
ENV FILE_OUTPUT=true EXECUTABLE_PATH=/usr/bin/chromium
