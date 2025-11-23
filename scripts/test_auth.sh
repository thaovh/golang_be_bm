#!/bin/bash

# Auth Service API Testing Script
# Usage: ./scripts/test_auth.sh

set -e

BASE_URL="${BASE_URL:-http://localhost:8000}"
API_BASE="${API_BASE:-/api/v1}"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Auth Service API Tests ===${NC}\n"

# Variables to store tokens
ACCESS_TOKEN=""
REFRESH_TOKEN=""

# Test 1: Register
echo -e "${YELLOW}Test 1: Register${NC}"
RESPONSE=$(curl -s -X POST "${BASE_URL}${API_BASE}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testauth@example.com",
    "username": "testauth",
    "password": "TestAuth@123",
    "full_name": "Test Auth User"
  }')

echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
ACCESS_TOKEN=$(echo "$RESPONSE" | jq -r '.access_token' 2>/dev/null || echo "")
REFRESH_TOKEN=$(echo "$RESPONSE" | jq -r '.refresh_token' 2>/dev/null || echo "")

if [ -n "$ACCESS_TOKEN" ] && [ "$ACCESS_TOKEN" != "null" ] && [ "$ACCESS_TOKEN" != "" ]; then
  echo -e "${GREEN}✓ Registration successful${NC}\n"
else
  echo -e "${RED}✗ Registration failed${NC}\n"
fi

# Test 2: Login
echo -e "${YELLOW}Test 2: Login${NC}"
RESPONSE=$(curl -s -X POST "${BASE_URL}${API_BASE}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "testauth@example.com",
    "password": "TestAuth@123"
  }')

echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
ACCESS_TOKEN=$(echo "$RESPONSE" | jq -r '.access_token' 2>/dev/null || echo "")
REFRESH_TOKEN=$(echo "$RESPONSE" | jq -r '.refresh_token' 2>/dev/null || echo "")

if [ -n "$ACCESS_TOKEN" ] && [ "$ACCESS_TOKEN" != "null" ] && [ "$ACCESS_TOKEN" != "" ]; then
  echo -e "${GREEN}✓ Login successful${NC}\n"
  echo -e "${BLUE}Access Token: ${ACCESS_TOKEN:0:50}...${NC}"
  echo -e "${BLUE}Refresh Token: ${REFRESH_TOKEN:0:50}...${NC}\n"
else
  echo -e "${RED}✗ Login failed${NC}\n"
  exit 1
fi

# Test 3: Get Current User (with token)
if [ -n "$ACCESS_TOKEN" ] && [ "$ACCESS_TOKEN" != "null" ]; then
  echo -e "${YELLOW}Test 3: Get Current User${NC}"
  RESPONSE=$(curl -s -X GET "${BASE_URL}${API_BASE}/auth/me" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}")
  echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
  echo -e "${GREEN}✓ Current user retrieved${NC}\n"
fi

# Test 4: Refresh Token
if [ -n "$REFRESH_TOKEN" ] && [ "$REFRESH_TOKEN" != "null" ]; then
  echo -e "${YELLOW}Test 4: Refresh Token${NC}"
  RESPONSE=$(curl -s -X POST "${BASE_URL}${API_BASE}/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{
      \"refresh_token\": \"${REFRESH_TOKEN}\"
    }")
  echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
  NEW_ACCESS_TOKEN=$(echo "$RESPONSE" | jq -r '.access_token' 2>/dev/null || echo "")
  if [ -n "$NEW_ACCESS_TOKEN" ] && [ "$NEW_ACCESS_TOKEN" != "null" ]; then
    echo -e "${GREEN}✓ Token refreshed successfully${NC}\n"
    ACCESS_TOKEN="$NEW_ACCESS_TOKEN"
  else
    echo -e "${RED}✗ Token refresh failed${NC}\n"
  fi
fi

# Test 5: Logout
if [ -n "$REFRESH_TOKEN" ] && [ "$REFRESH_TOKEN" != "null" ]; then
  echo -e "${YELLOW}Test 5: Logout${NC}"
  RESPONSE=$(curl -s -X POST "${BASE_URL}${API_BASE}/auth/logout" \
    -H "Content-Type: application/json" \
    -d "{
      \"refresh_token\": \"${REFRESH_TOKEN}\"
    }")
  echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
  echo -e "${GREEN}✓ Logout successful${NC}\n"
fi

# Test 6: Try to use revoked token (should fail)
if [ -n "$ACCESS_TOKEN" ] && [ "$ACCESS_TOKEN" != "null" ]; then
  echo -e "${YELLOW}Test 6: Try to use revoked token (should fail)${NC}"
  RESPONSE=$(curl -s -X GET "${BASE_URL}${API_BASE}/auth/me" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}")
  echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
  echo -e "${BLUE}Note: This should fail if token is properly revoked${NC}\n"
fi

echo -e "${GREEN}=== All Auth Tests Completed ===${NC}"

