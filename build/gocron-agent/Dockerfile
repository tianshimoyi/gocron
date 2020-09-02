FROM golang:1.14-alpine AS builder
RUN apk update && mkdir /app
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 go build -o  dist/gocron-agent -a -ldflags "-w -s" ./cmd/gocron-agent

FROM alpine:3.10

ARG REPO_URL
ARG BRANCH
ARG COMMIT_REF
LABEL repo-url=${REPO_URL}
LABEL branch=${BRANCH}
LABEL commit-ref=${COMMIT_REF}

RUN apk update \
  && apk upgrade \
  && apk add --no-cache \
  ca-certificates tzdata bash \
  && update-ca-certificates 2>/dev/null || true \
  && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
  && mkdir /app
WORKDIR /app
COPY --from=builder /app/dist/gocron-agent /app/gocron-agent

EXPOSE 8080

CMD ["/app/gocron-agent"]