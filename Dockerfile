### BUILDER STAGE
#
#
FROM golang:1.27-alpine3.23 AS builder

# Do not remove `git` here, it is required for getting runner version when executing `make build`
RUN apk add --no-cache make git

ARG GOPROXY
ENV GOPROXY=${GOPROXY:-}

COPY . /opt/src/runner
WORKDIR /opt/src/runner

RUN make clean && make build

### BASIC VARIANT
#
#
FROM alpine:3.24

ARG VERSION=dev

LABEL org.opencontainers.image.source="https://gitea.com/gitea/runner"
LABEL org.opencontainers.image.version="${VERSION}"

RUN apk add --no-cache bash git tzdata

COPY --from=builder /opt/src/runner/gitea-runner /usr/local/bin/gitea-runner
COPY scripts/run.sh /usr/local/bin/run.sh

VOLUME /data

ENTRYPOINT ["run.sh"]
