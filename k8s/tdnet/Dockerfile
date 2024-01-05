# localhost:5000/tdnet:latest
FROM alpine:latest as clone
RUN apk add --update git
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/browser-operator.git --depth 1

FROM node:18-alpine as build
COPY --from=clone /app/browser-operator/tdnet /app/tdnet
WORKDIR /app/tdnet
RUN apk add --update curl chromium
RUN npm install && npm install -g typescript && tsc
ENV FILE_OUTPUT=true
ENV EXECUTABLE_PATH=/usr/bin/chromium
