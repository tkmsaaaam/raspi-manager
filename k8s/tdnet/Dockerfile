# localhost:5000/tdnet:latest
FROM alpine:latest as clone
RUN apk add --update git
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/browser-operator.git --depth 1

FROM node:21-alpine
COPY --from=clone /app/browser-operator/tdnet /app/tdnet
WORKDIR /app/tdnet
RUN apk add --update curl
ENV FILE_OUTPUT=true EXECUTABLE_PATH=/usr/bin/chromium
