# clamavreceiver:latest
FROM golang:1.22-bookworm as builder
WORKDIR /app
RUN git clone https://github.com/open-telemetry/opentelemetry-collector.git --depth=1 && git clone https://github.com/tkmsaaaam/raspi-manager.git --depth=1
WORKDIR /app/opentelemetry-collector/cmd/mdatagen
RUN go install .
WORKDIR /app/raspi-manager/receiver/clamavreceiver
RUN go generate ./...
WORKDIR /app/raspi-manager/receiver
RUN go install go.opentelemetry.io/collector/cmd/builder@v0.107.0 && builder --config=otelcol-builder.yaml

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/raspi-manager/receiver/clamav /app/clamav
CMD ["/bin/bash", "-c", "/app/clamav --config=/config/config.yaml"]
