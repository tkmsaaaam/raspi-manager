# localhost:5000/count-commits-js:latest
FROM alpine:latest as clone
RUN apk update && apk add git
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/raspi-manager.git --depth 1
RUN git clone https://github.com/tkmsaaaam/count-commits-js.git --depth 1

FROM golang:1.21-alpine as build-main
COPY --from=clone /app/count-commits-js /app/count-commits-js
WORKDIR /app/count-commits-js
RUN go mod tidy
RUN go build

FROM golang:1.21-alpine as build-jwt
COPY --from=clone /app/raspi-manager/k8s/countCommitsJs/jwt /app/jwt
WORKDIR /app/jwt
RUN go mod tidy
RUN go build

FROM alpine:latest
COPY --from=build-jwt /app/jwt/jwt /app/jwt
COPY --from=build-main /app/count-commits-js/count-commits-js /app/count-commits-js
