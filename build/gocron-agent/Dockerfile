FROM golang:1.14-alpine AS builder
RUN apk update && mkdir /app
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 go build -o  dist/gocron-agent -a -ldflags "-w -s" ./cmd/gocron-agent

FROM alpine:3.10
RUN apk update \
  && apk upgrade \
  && apk add --no-cache \
  ca-certificates \
  && update-ca-certificates 2>/dev/null || true \
  && mkdir /app
WORKDIR /app
COPY --from=builder /app/dist/gocron-agent /app/gocron-agent

EXPOSE 8080

CMD ["/app/gocron-agent"]