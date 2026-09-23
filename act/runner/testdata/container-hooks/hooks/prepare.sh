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

ACTIONS_COUNT=$(jq '.runner.actions // [] | length' "$INPUT")
i=0
while [ "$i" -lt "$ACTIONS_COUNT" ]; do
    HASH=$(jq -r ".runner.actions[$i].hash" "$INPUT")
    SRC=$(jq -r ".runner.actions[$i].localPath" "$INPUT")
    DST="$ACT_PATH/actions/$HASH"

    rm -rf "$DST"
    mkdir -p "$(dirname "$DST")"
    cp -r "$SRC" "$DST"

    i=$((i + 1))
done

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