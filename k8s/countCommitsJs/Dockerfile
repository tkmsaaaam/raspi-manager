# localhost:5000/count-commits-js
FROM alpine:3.12 as clone
RUN apk update && apk add git
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/raspi-manager.git
RUN git clone https://github.com/tkmsaaaam/count-commits-js.git

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

FROM alpine:3.12
COPY --from=build-jwt /app/jwt/jwt /app/jwt
COPY --from=build-main /app/count-commits-js/count-commits-js /app/count-commits-js
