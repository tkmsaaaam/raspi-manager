# localhost:5000/check-login
FROM alpine:latest as clone
RUN apk update && apk add git
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/raspi-manager.git --depth 1

FROM mcr.microsoft.com/dotnet/sdk:8.0-alpine
COPY --from=clone /app/raspi-manager/k8s/checkLogin/AuthlogParser /app/AuthlogParser
WORKDIR /app/AuthlogParser
