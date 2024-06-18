# clamavreceiver:latest
FROM golang:1.22-bookworm as builder
WORKDIR /app
COPY . .
RUN go install go.opentelemetry.io/collector/cmd/builder@latest && builder --config=otelcol-builder.yaml

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/clamav /app/clamav
CMD ["/bin/bash", "-c", "/app/clamav --config=/config/config.yaml"]
