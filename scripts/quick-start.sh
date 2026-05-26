#!/usr/bin/env bash
set -euo pipefail

# Quick start script for the project
# go run ./cmd/migrate create create_whitelist_ips_table
# go run ./cmd/migrate up
# swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal


SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "${PROJECT_ROOT}"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo ""
echo "Starting database services (db + pgadmin)..."
echo ""

docker compose up -d db pgadmin

echo ""
echo "Waiting for database to be healthy..."
echo ""

MAX_RETRIES=20
RETRY_DELAY=3
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  if docker compose ps db | grep -q "healthy"; then
    echo -e "${GREEN}✅ Database is healthy${NC}"
    break
  fi

  RETRY_COUNT=$((RETRY_COUNT + 1))
  if [ $RETRY_COUNT -lt $MAX_RETRIES ]; then
    echo -e "${YELLOW}⏳ Waiting for db health... (${RETRY_COUNT}/${MAX_RETRIES})${NC}"
    sleep $RETRY_DELAY
  else
    echo -e "${RED}❌ Database did not become healthy in time${NC}"
    echo "Check logs with: docker compose logs db"
    exit 1
  fi
done

echo ""
echo "Running database migrations..."
echo ""

MAX_MIGRATION_RETRIES=3
MIGRATION_RETRY_DELAY=3
MIGRATION_RETRY_COUNT=0
MIGRATION_SUCCESS=false

while [ $MIGRATION_RETRY_COUNT -lt $MAX_MIGRATION_RETRIES ]; do
  if go run "${PROJECT_ROOT}/cmd/migrate/main.go" up; then
    MIGRATION_SUCCESS=true
    break
  else
    MIGRATION_RETRY_COUNT=$((MIGRATION_RETRY_COUNT + 1))
    if [ $MIGRATION_RETRY_COUNT -lt $MAX_MIGRATION_RETRIES ]; then
      echo -e "${YELLOW}⚠️ Migration attempt ${MIGRATION_RETRY_COUNT} failed, retrying in ${MIGRATION_RETRY_DELAY}s...${NC}"
      sleep $MIGRATION_RETRY_DELAY
    fi
  fi
done

if [ "$MIGRATION_SUCCESS" = true ]; then
  echo ""
  echo -e "${GREEN}✅ Migrations completed successfully${NC}"
  echo -e "${GREEN}✅ Quick start done${NC}"
  echo ""
  echo "Run app locally with:"
  echo "  APP_PORT=8081 go run ./cmd/server"
  echo ""
else
  echo ""
  echo -e "${RED}❌ Failed to run migrations after ${MAX_MIGRATION_RETRIES} attempts${NC}"
  echo "Check logs with: docker compose logs db"
  exit 1
fi