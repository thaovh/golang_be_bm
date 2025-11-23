#!/bin/bash

# Database migration script
# Usage: ./scripts/migrate.sh [migration_file]

set -e

# Database connection details
DB_HOST="${DB_HOST:-127.0.0.1}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-bm_staff}"
MIGRATIONS_DIR="${MIGRATIONS_DIR:-./migrations}"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Running database migrations...${NC}"
echo "Database: ${DB_NAME}@${DB_HOST}:${DB_PORT}"
echo "User: ${DB_USER}"
echo ""

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo -e "${RED}Error: psql command not found. Please install PostgreSQL client.${NC}"
    exit 1
fi

# If specific migration file is provided, run only that
if [ -n "$1" ]; then
    MIGRATION_FILE="$1"
    if [ ! -f "$MIGRATION_FILE" ]; then
        echo -e "${RED}Error: Migration file not found: $MIGRATION_FILE${NC}"
        exit 1
    fi
    echo -e "${YELLOW}Running migration: $MIGRATION_FILE${NC}"
    PGPASSWORD="${DB_PASSWORD}" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$MIGRATION_FILE"
    echo -e "${GREEN}Migration completed successfully!${NC}"
    exit 0
fi

# Run all migrations in order
if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo -e "${RED}Error: Migrations directory not found: $MIGRATIONS_DIR${NC}"
    exit 1
fi

# Find all SQL files and sort them
MIGRATION_FILES=$(find "$MIGRATIONS_DIR" -name "*.sql" | sort)

if [ -z "$MIGRATION_FILES" ]; then
    echo -e "${YELLOW}No migration files found in $MIGRATIONS_DIR${NC}"
    exit 0
fi

# Run each migration
for file in $MIGRATION_FILES; do
    echo -e "${YELLOW}Running migration: $(basename $file)${NC}"
    PGPASSWORD="${DB_PASSWORD}" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$file"
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ $(basename $file) completed${NC}"
    else
        echo -e "${RED}✗ $(basename $file) failed${NC}"
        exit 1
    fi
done

echo ""
echo -e "${GREEN}All migrations completed successfully!${NC}"

