# localhost:5000/unitedcinemas:latest
FROM node:21-alpine
RUN apk add --update --no-cache git
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/browser-operator.git --depth 1
WORKDIR /app/browser-operator/unitedcinemas
