#!/bin/sh
set -e

echo "prepare hook: start"

INPUT="$RUNNER_CONTAINER_HOOK_INPUT_FILE"
RESPONSE="$RUNNER_CONTAINER_HOOK_RESPONSE_FILE"

if [ ! -f "$INPUT" ]; then
    echo "prepare hook: $INPUT not found" >&2
    exit 1
fi

cat "$INPUT"

WORKDIR=$(jq -r '.runner.workdir' "$INPUT")

ACT_PATH="$RUNNER_CONTAINER_HOOK_STATE_DIR/act"
mkdir -p "$ACT_PATH/workflow"

mkdir -p "$WORKDIR"

jq -r '.runner.eventJSON' "$INPUT" > "$ACT_PATH/workflow/event.json"
: > "$ACT_PATH/workflow/envs.txt"

cat > "$RESPONSE" <<EOF
{
  "id": "test-pod-123",
  "actPath": "$ACT_PATH",
  "environments": {
    "CUSTOM_VAR_FROM_PREPARE": "hello-from-prepare"
  }
}
EOF

echo "prepare hook: actPath=$ACT_PATH"
echo "prepare hook: done"