# manage-slack-daily:latest
FROM debian:bookworm-slim as clone
RUN apt-get update -y && apt-get install --no-install-recommends -y git ca-certificates
WORKDIR /app
RUN git clone https://github.com/tkmsaaaam/manage-slack.git --depth 1

FROM golang:1.22-bookworm as build
COPY --from=clone /app/manage-slack /app/manage-slack
WORKDIR /app/manage-slack
RUN go mod tidy && go build daily.go

FROM debian:bookworm-slim
RUN apt-get update -y && apt-get install --no-install-recommends -y ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*
COPY --from=build /app/manage-slack/daily /app/daily
