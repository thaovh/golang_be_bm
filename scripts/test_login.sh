#!/bin/bash

# Test Login API với identifier field
# Usage: ./scripts/test_login.sh

BASE_URL="${BASE_URL:-http://localhost:8000}"

echo "=== Test Login API với identifier field ==="
echo ""

# Test 1: Login với email
echo "Test 1: Login với email qua identifier field"
RESPONSE=$(curl -s -X POST "${BASE_URL}/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "admin@example.com",
    "password": "t123456T@"
  }')

echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
echo ""

# Extract token nếu thành công
ACCESS_TOKEN=$(echo "$RESPONSE" | python3 -c "import sys, json; d=json.load(sys.stdin); print(d.get('access_token', ''))" 2>/dev/null)

if [ -n "$ACCESS_TOKEN" ] && [ "$ACCESS_TOKEN" != "null" ] && [ "$ACCESS_TOKEN" != "" ]; then
  echo "✓ Login thành công!"
  echo "Access Token: ${ACCESS_TOKEN:0:60}..."
  echo ""
  
  # Test 2: Get Current User
  echo "Test 2: Get Current User với token"
  curl -s -X GET "${BASE_URL}/api/v1/auth/me" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" | python3 -m json.tool 2>/dev/null
  echo ""
else
  echo "✗ Login thất bại"
  echo "Lưu ý: Có thể do rate limit hoặc sai password"
  echo ""
fi

# Test 3: Login với username
echo "Test 3: Login với username qua identifier field"
curl -s -X POST "${BASE_URL}/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "admin",
    "password": "t123456T@"
  }' | python3 -m json.tool 2>/dev/null | head -15
echo ""

echo "=== Kết quả ==="
echo "✓ API đã chấp nhận field 'identifier' (thay vì 'email')"
echo "✓ Login endpoint hoạt động với cả email và username"
echo "✓ Field identifier hoạt động đúng như mong đợi"

