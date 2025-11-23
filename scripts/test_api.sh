#!/bin/bash

# API Testing Script
# Usage: ./scripts/test_api.sh

set -e

BASE_URL="${BASE_URL:-http://localhost:8000}"
API_BASE="${API_BASE:-/api/v1}"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== User Service API Tests ===${NC}\n"

# Test 1: Create User
echo -e "${YELLOW}Test 1: Create User${NC}"
RESPONSE=$(curl -s -X POST "${BASE_URL}${API_BASE}/users" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "Test@1234",
    "full_name": "Test User",
    "gender": "male",
    "role": "user"
  }')

echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
USER_ID=$(echo "$RESPONSE" | jq -r '.user.id' 2>/dev/null || echo "")

if [ -n "$USER_ID" ] && [ "$USER_ID" != "null" ]; then
  echo -e "${GREEN}✓ User created successfully: $USER_ID${NC}\n"
else
  echo -e "${RED}✗ Failed to create user${NC}\n"
fi

# Test 2: Get User by ID
if [ -n "$USER_ID" ] && [ "$USER_ID" != "null" ]; then
  echo -e "${YELLOW}Test 2: Get User by ID${NC}"
  RESPONSE=$(curl -s -X GET "${BASE_URL}${API_BASE}/users/${USER_ID}")
  echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
  echo -e "${GREEN}✓ User retrieved${NC}\n"
fi

# Test 3: Get User by Email
echo -e "${YELLOW}Test 3: Get User by Email${NC}"
RESPONSE=$(curl -s -X GET "${BASE_URL}${API_BASE}/users/email/test@example.com")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo -e "${GREEN}✓ User retrieved by email${NC}\n"

# Test 4: Get User by Username
echo -e "${YELLOW}Test 4: Get User by Username${NC}"
RESPONSE=$(curl -s -X GET "${BASE_URL}${API_BASE}/users/username/testuser")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo -e "${GREEN}✓ User retrieved by username${NC}\n"

# Test 5: List Users
echo -e "${YELLOW}Test 5: List Users${NC}"
RESPONSE=$(curl -s -X GET "${BASE_URL}${API_BASE}/users?page=1&page_size=10")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo -e "${GREEN}✓ Users listed${NC}\n"

# Test 6: Update User
if [ -n "$USER_ID" ] && [ "$USER_ID" != "null" ]; then
  echo -e "${YELLOW}Test 6: Update User${NC}"
  RESPONSE=$(curl -s -X PUT "${BASE_URL}${API_BASE}/users/${USER_ID}" \
    -H "Content-Type: application/json" \
    -d '{
      "full_name": "Updated Test User",
      "gender": "female"
    }')
  echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
  echo -e "${GREEN}✓ User updated${NC}\n"
fi

# Test 7: Change Password
if [ -n "$USER_ID" ] && [ "$USER_ID" != "null" ]; then
  echo -e "${YELLOW}Test 7: Change Password${NC}"
  RESPONSE=$(curl -s -X POST "${BASE_URL}${API_BASE}/users/${USER_ID}/change-password" \
    -H "Content-Type: application/json" \
    -d '{
      "new_password": "NewPass@1234"
    }')
  echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
  echo -e "${GREEN}✓ Password changed${NC}\n"
fi

# Test 8: Validation Tests
echo -e "${YELLOW}Test 8: Validation Tests${NC}"

# Invalid email
echo -e "${BLUE}Testing invalid email...${NC}"
RESPONSE=$(curl -s -X POST "${BASE_URL}${API_BASE}/users" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid-email",
    "username": "testuser2",
    "password": "Test@1234"
  }')
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Weak password
echo -e "${BLUE}Testing weak password...${NC}"
RESPONSE=$(curl -s -X POST "${BASE_URL}${API_BASE}/users" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test2@example.com",
    "username": "testuser2",
    "password": "weak"
  }')
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Invalid username
echo -e "${BLUE}Testing invalid username...${NC}"
RESPONSE=$(curl -s -X POST "${BASE_URL}${API_BASE}/users" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test3@example.com",
    "username": "ab",
    "password": "Test@1234"
  }')
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 9: Delete User (optional - comment out if you want to keep test data)
# if [ -n "$USER_ID" ] && [ "$USER_ID" != "null" ]; then
#   echo -e "${YELLOW}Test 9: Delete User${NC}"
#   RESPONSE=$(curl -s -X DELETE "${BASE_URL}${API_BASE}/users/${USER_ID}")
#   echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
#   echo -e "${GREEN}✓ User deleted${NC}\n"
# fi

echo -e "${GREEN}=== All Tests Completed ===${NC}"

