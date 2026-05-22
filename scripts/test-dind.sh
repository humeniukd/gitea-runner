#!/usr/bin/env bash
# Generic docker-in-docker test harness.
#
# Builds a dind image variant from the repo Dockerfile, starts its docker daemon over a
# local TCP port, and runs a Go test command against that daemon via DOCKER_HOST. This
# validates the actual docker version and behaviour shipped in the dind image, so any
# daemon-level regression surfaces here (e.g. the "docker cp" break in gitea/runner#981).
# It is deliberately generic: point it at any package/test to exercise the dind daemon.
#
# Usage: scripts/test-dind.sh [target] [-- go-test-args...]
#   target:        dind (default) or dind-rootless
#   go-test-args:  passed verbatim to `go test`
#                  (default: -race ./act/container/, the daemon-facing package)
#
# Env:
#   DIND_TEST_PORT  host port for the daemon (default 32375)
#   DIND_TEST_IMAGE skip the build and use this prebuilt image instead
set -euo pipefail

target="dind"
case "${1:-}" in
  dind|dind-rootless) target="$1"; shift ;;
esac
[ "${1:-}" = "--" ] && shift
test_args=("$@")
[ ${#test_args[@]} -eq 0 ] && test_args=(-race ./act/container/)

port="${DIND_TEST_PORT:-32375}"
name="gitea-runner-dind-test-$$"
image="${DIND_TEST_IMAGE:-gitea-runner-${target}:dind-test}"

cleanup() { docker rm -f "$name" >/dev/null 2>&1 || true; }
trap cleanup EXIT

if [ -z "${DIND_TEST_IMAGE:-}" ]; then
  echo "==> Building ${target} image"
  docker build --target "$target" -t "$image" .
fi

# Override the image entrypoint (s6) and run only dockerd, exposed over insecure TCP.
# We are testing the daemon the image ships, not the runner supervision tree.
#
# How the test process reaches the daemon depends on where it runs:
#   - plain host: publish 2375 on loopback and connect to 127.0.0.1.
#   - inside a container (CI), the daemon is a sibling container, so its published port is on
#     the host, not our loopback; instead attach it to our own network and reach it by name.
self_container=""
if [ -f /.dockerenv ]; then
  self_container="$(cat /proc/sys/kernel/hostname 2>/dev/null || cat /etc/hostname)"
fi
self_network=""
if [ -n "$self_container" ]; then
  self_network="$(docker inspect -f '{{range $k,$v := .NetworkSettings.Networks}}{{$k}}{{"\n"}}{{end}}' "$self_container" 2>/dev/null | head -1)"
fi

if [ -n "$self_network" ]; then
  echo "==> Starting ${target} daemon on network ${self_network} (reached as ${name}:2375)"
  docker run -d --privileged --name "$name" --network "$self_network" \
    -e DOCKER_TLS_CERTDIR= \
    --entrypoint dockerd-entrypoint.sh \
    "$image" --host=tcp://0.0.0.0:2375 >/dev/null
  export DOCKER_HOST="tcp://${name}:2375"
else
  echo "==> Starting ${target} daemon on tcp://127.0.0.1:${port}"
  docker run -d --privileged --name "$name" \
    -e DOCKER_TLS_CERTDIR= \
    -p "127.0.0.1:${port}:2375" \
    --entrypoint dockerd-entrypoint.sh \
    "$image" --host=tcp://0.0.0.0:2375 >/dev/null
  export DOCKER_HOST="tcp://127.0.0.1:${port}"
fi

echo "==> Waiting for daemon"
for _ in $(seq 1 60); do
  docker version >/dev/null 2>&1 && break
  sleep 1
done
docker version --format 'server docker {{.Server.Version}}'

echo "==> Running tests against dind daemon"
go test "${test_args[@]}"
