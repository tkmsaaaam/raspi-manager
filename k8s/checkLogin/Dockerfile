# localhost:5000/check-login
FROM alpine:3.12 as clone
RUN apk update && apk add git
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/raspi-manager.git

FROM mcr.microsoft.com/dotnet/sdk:8.0 as build
COPY --from=clone /app/raspi-manager/k8s/checkLogin/AuthlogParser /app/AuthlogParser
WORKDIR /app/AuthlogParser
