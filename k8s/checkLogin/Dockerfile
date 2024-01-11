# localhost:5000/check-login:latest
FROM alpine:latest as clone
RUN apk add --update git
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/raspi-manager.git --depth 1

FROM mcr.microsoft.com/dotnet/sdk:8.0-alpine as build
COPY --from=clone /app/raspi-manager/k8s/checkLogin/AuthlogParser /app/AuthlogParser
WORKDIR /app/AuthlogParser
RUN dotnet publish -c Release --self-contained true /p:PublishSingleFile=true

FROM alpine:latest
COPY --from=build /app/AuthlogParser/bin/Release/net8.0/linux-musl-arm/publish /app
RUN apk add --update libgcc libstdc++ icu-libs
WORKDIR /app
