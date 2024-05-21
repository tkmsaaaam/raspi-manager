# unitedcinemas:latest
FROM debian:bookworm-slim as build
WORKDIR /app
RUN apt-get update -y && apt-get install --no-install-recommends -y git ca-certificates
RUN git clone https://github.com/tkmsaaaam/browser-operator.git --depth 1
WORKDIR /app/browser-operator/unitedcinemas

FROM node:21-bookworm-slim
COPY --from=build /app/browser-operator/unitedcinemas /app/unitedcinemas
RUN apt-get update -y  && apt-get install --no-install-recommends -y curl ca-certificates  && apt-get clean && rm -rf /var/lib/apt/lists/*
WORKDIR /app/unitedcinemas
