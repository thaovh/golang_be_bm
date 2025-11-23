#!/bin/bash

# Script to create admin user
# Usage: ./scripts/create_admin.sh [email] [username] [password]

set -e

BASE_URL="${BASE_URL:-http://localhost:8000}"
API_BASE="${API_BASE:-/api/v1}"

# Default values
EMAIL="${1:-admin@example.com}"
USERNAME="${2:-admin}"
PASSWORD="${3:-Admin@123}"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Creating admin user...${NC}"
echo "Email: $EMAIL"
echo "Username: $USERNAME"
echo ""

# Check if service is running
if ! curl -s -f "${BASE_URL}/api/v1/users?page=1&page_size=1" > /dev/null 2>&1; then
    echo -e "${RED}Error: Service is not running at ${BASE_URL}${NC}"
    echo "Please start the service first:"
    echo "  go run cmd/server/main.go cmd/server/wire_gen.go"
    exit 1
fi

# Create admin user
RESPONSE=$(curl -s -X POST "${BASE_URL}${API_BASE}/users" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"${EMAIL}\",
    \"username\": \"${USERNAME}\",
    \"password\": \"${PASSWORD}\",
    \"full_name\": \"Administrator\",
    \"role\": \"admin\"
  }")

# Check if jq is available
if command -v jq &> /dev/null; then
    USER_ID=$(echo "$RESPONSE" | jq -r '.user.id' 2>/dev/null)
    ERROR=$(echo "$RESPONSE" | jq -r '.message // .error // empty' 2>/dev/null)
    
    if [ -n "$USER_ID" ] && [ "$USER_ID" != "null" ]; then
        echo -e "${GREEN}✓ Admin user created successfully!${NC}"
        echo ""
        echo "$RESPONSE" | jq '.'
        echo ""
        echo -e "${YELLOW}User ID: ${USER_ID}${NC}"
        echo -e "${YELLOW}Email: ${EMAIL}${NC}"
        echo -e "${YELLOW}Username: ${USERNAME}${NC}"
        echo -e "${YELLOW}Password: ${PASSWORD}${NC}"
    else
        echo -e "${RED}✗ Failed to create admin user${NC}"
        echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
        if [ -n "$ERROR" ]; then
            echo -e "${RED}Error: $ERROR${NC}"
        fi
        exit 1
    fi
else
    echo "$RESPONSE"
    if echo "$RESPONSE" | grep -q '"id"'; then
        echo -e "${GREEN}✓ Admin user created successfully!${NC}"
    else
        echo -e "${RED}✗ Failed to create admin user${NC}"
        exit 1
    fi
fi

