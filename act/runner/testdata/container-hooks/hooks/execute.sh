#!/bin/sh
set -e

echo "execute hook: start"

INPUT="$RUNNER_CONTAINER_HOOK_INPUT_FILE"
RESPONSE="$RUNNER_CONTAINER_HOOK_RESPONSE_FILE"

if [ ! -f "$INPUT" ]; then
    echo "execute hook: $INPUT not found" >&2
    exit 1
fi

cat "$INPUT"

KIND=$(jq -r '.kind' "$INPUT")
WORKDIR=$(jq -r '.workingDir // ""' "$INPUT")

if [ -n "$WORKDIR" ] && [ -d "$WORKDIR" ]; then
    cd "$WORKDIR"
fi

# Export step env from JSON (INPUT_*, workflow env, etc.). Names can contain
# hyphens, so use the env command rather than shell export.
ENV_ARGS=$(jq -r '.env // {} | to_entries | map((.key + "=" + .value) | @sh) | .[]' "$INPUT" 2>/dev/null | tr '\n' ' ')

# Create file-command files before the command runs.
FILE_COMMANDS=$(jq -r '.fileCommands[]? // empty' "$INPUT")
if [ -n "$FILE_COMMANDS" ]; then
    echo "$FILE_COMMANDS" | while read -r f; do
        [ -z "$f" ] && continue
        mkdir -p "$(dirname "$f")"
        : > "$f"
    done
fi

EXIT_CODE=0

case "$KIND" in
    script)
        SCRIPT=$(jq -r '.script' "$INPUT")
        printf '%s' "$SCRIPT" | eval "env $ENV_ARGS sh -e -" || EXIT_CODE=$?
        ;;
    command)
        CMD_SH=$(jq -r '.command[] | @sh' "$INPUT" | tr '\n' ' ')
        eval "env $ENV_ARGS $CMD_SH" || EXIT_CODE=$?
        ;;
    container-step)
        ENTRY_SH=$(jq -r '.entrypoint[]? | @sh' "$INPUT" | tr '\n' ' ')
        ARGS_SH=$(jq -r '.args[]? | @sh' "$INPUT" | tr '\n' ' ')
        if [ -z "$ENTRY_SH$ARGS_SH" ]; then
            EXIT_CODE=1
        else
            eval "env $ENV_ARGS $ENTRY_SH $ARGS_SH" || EXIT_CODE=$?
        fi
        ;;
    *)
        echo "execute hook: unknown kind '$KIND'" >&2
        exit 1
        ;;
esac

# Collect file-command contents into the response.
if [ -n "$RESPONSE" ]; then
    {
        printf '{"exitCode":%d,"fileCommands":{' "$EXIT_CODE"
        first=1
        if [ -n "$FILE_COMMANDS" ]; then
            echo "$FILE_COMMANDS" | while read -r f; do
                [ -z "$f" ] && continue
                if [ -f "$f" ]; then
                    if [ $first -eq 0 ]; then printf ','; fi
                    first=0
                    content=$(sed 's/\\/\\\\/g; s/"/\\"/g' "$f" | awk '{printf "%s\\n", $0}')
                    printf '"%s":"%s"' "$f" "$content"
                fi
            done
        fi
        printf '}}'
    } > "$RESPONSE"
fi

echo "execute hook: done"
exit 0