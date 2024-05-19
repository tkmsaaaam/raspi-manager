# check-login:latest

# FROM debian:bookworm-slim
# COPY AuthlogParser/bin/Release/net8.0/linux-arm64/publish/ /app/
# RUN apt-get update -y  && apt-get install --no-install-recommends -y curl tzdata jq ca-certificates libc6 libgcc-s1 libgssapi-krb5-2 libicu72 libstdc++6 libssl1.1 && apt-get clean && rm -rf /var/lib/apt/lists/*
# WORKDIR /app
# CMD ["/bin/sh", "/app/AuthlogParser"]
FROM alpine:latest
COPY AuthlogParser/bin/Release/net8.0/linux-musl-arm/publish/ /app/
RUN apk add --update --no-cache libgcc libssl3 libstdc++ zlib icu-libs tzdata curl jq
WORKDIR /app
