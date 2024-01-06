# localhost:5000/rakuten-sec:latest
FROM alpine:latest as clone
RUN apk add --update git
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/browser-operator.git --depth 1

FROM node:21-alpine as build
COPY --from=clone /app/browser-operator/rakuten-sec /app/rakuten-sec
WORKDIR /app/rakuten-sec
RUN apk add --update curl chromium
ENV FILE_OUTPUT=true EXECUTABLE_PATH=/usr/bin/chromium
