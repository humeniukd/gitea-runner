#!/bin/sh
set -e

DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"

case "${RUNNER_CONTAINER_HOOK_TYPE:-}" in
    prepare) exec "$DIR/prepare.sh" ;;
    execute) exec "$DIR/execute.sh" ;;
    cleanup) exec "$DIR/cleanup.sh" ;;
    *)
        echo "main.sh: unknown RUNNER_CONTAINER_HOOK_TYPE='${RUNNER_CONTAINER_HOOK_TYPE:-}'" >&2
        exit 1
        ;;
esac