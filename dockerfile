FROM golang:1.23.2-alpine3.20 AS builder

ARG DOCKER_WORKDIR=/app

WORKDIR ${DOCKER_WORKDIR}

COPY . .

RUN apk add bash

RUN go mod tidy

RUN bash -c "go build -o ./bin/cybertask ./cmd/main.go"

FROM alpine:3.20 AS runner

ARG DOCKER_WORKDIR=/app

WORKDIR ${DOCKER_WORKDIR}

COPY --from=builder ${DOCKER_WORKDIR}/bin/cybertask ./bin/

ENV WORKDIR ${DOCKER_WORKDIR}

ENV CONFIG_FILE ${WORKDIR}/config.yaml

CMD ./bin/cybertask 