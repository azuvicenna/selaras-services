#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${1:-${BFF_URL:-http://localhost:8080}}"
VERBOSE="${2:-}"

echo -e "\033[1;36m==========================================================\033[0m"
echo -e "\033[1;36m  Selaras SIMRS - Executing E2E Patient Journey (Bash)    \033[0m"
echo -e "\033[1;36m==========================================================\033[0m"
echo -e "\033[1;33mTarget Gateway: ${BASE_URL}\033[0m"

FLAGS=("-base-url" "${BASE_URL}")
if [[ -n "${VERBOSE}" ]]; then
    FLAGS+=("-verbose")
fi

go run ./tests/e2e/cmd/main.go "${FLAGS[@]}"
