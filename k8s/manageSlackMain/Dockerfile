# manage-slack-main:latest
FROM debian:bookworm-slim AS clone
RUN apt-get update -y && apt-get install --no-install-recommends -y git ca-certificates
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/manage-slack.git --depth 1

FROM golang:1.22-bookworm AS build
COPY --from=clone /app/manage-slack /app/manage-slack
WORKDIR /app/manage-slack
RUN go mod tidy && go build main.go

FROM debian:bookworm-slim
RUN apt-get update -y && apt-get install --no-install-recommends -y ca-certificates
COPY --from=build /app/manage-slack/main /app/main
