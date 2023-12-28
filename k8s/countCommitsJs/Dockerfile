FROM golang:1.16-alpine3.12 as build
RUN apk update && apk add git
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/raspi-manager.git
RUN git clone https://github.com/tkmsaaaam/count-commits-js.git
WORKDIR /app/count-commits-js
RUN go get
RUN go build

FROM alpine:3.12
COPY --from=build /app/raspi-manager/countCommitsJs/jwt /app/jwt
COPY --from=build /app/count-commits-js/count-commits-js /app/count-commits-js
RUN apk update && apk add npm
