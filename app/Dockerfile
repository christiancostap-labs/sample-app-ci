FROM golang:1.27.1-alpine3.24 AS builder

RUN apk add --no-cache git make gcc musl-dev


RUN go install go.opentelemetry.io/otelc/tool/cmd/otelc@latest

ENV OTEL_GO_ENABLED_INSTRUMENTATIONS=nethttp,runtimemetrics

WORKDIR /app

COPY . .

RUN go mod tidy
RUN otelc go build -a -o app .
    
FROM alpine:3.24

COPY --from=builder /app .

ENTRYPOINT ["/app"]