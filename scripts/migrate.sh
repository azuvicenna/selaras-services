#!/usr/bin/env bash
set -euo pipefail

ACTION="${1:-up}"
SERVICE="${2:-all}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"

echo "=========================================================="
echo " SELARAS HEALTHCARE ECOSYSTEM - UNIFIED MIGRATION RUNNER  "
echo " Action: ${ACTION} | Service: ${SERVICE} | Host: ${DB_HOST}:${DB_PORT}"
echo "=========================================================="

go run ./scripts/migrator -action="${ACTION}" -service="${SERVICE}" -db-host="${DB_HOST}" -db-port="${DB_PORT}"
