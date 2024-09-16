# count-commits-js:latest
FROM debian:bookworm-slim AS clone
RUN apt-get update -y && apt-get install --no-install-recommends -y git ca-certificates
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/raspi-manager.git --depth 1 && git clone https://github.com/tkmsaaaam/count-commits-js.git --depth 1

FROM golang:1.22-bookworm AS build-main
COPY --from=clone /app/count-commits-js /app/count-commits-js
WORKDIR /app/count-commits-js
RUN go mod tidy && go build

FROM golang:1.22-bookworm AS build-jwt
COPY --from=clone /app/raspi-manager/k8s/countCommitsJs/jwt /app/jwt
WORKDIR /app/jwt
RUN go mod tidy && go build

FROM debian:bookworm-slim
COPY --from=build-jwt /app/jwt/jwt /app/jwt
COPY --from=build-main /app/count-commits-js/count-commits-js /app/count-commits-js
RUN apt-get update -y && apt-get install --no-install-recommends -y ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*
