#!/bin/bash

# Pretty-print JSON output
function pretty_json {
  if command -v jq &>/dev/null; then
    echo "$1" | jq .
  else
    echo "$1" | sed 's/,/,\n/g' | sed 's/{/{\n/g' | sed 's/}/\n}/g'
  fi
}

# API base URL
API_URL="http://localhost:8080/api/v1"
echo "🚀 Testing Banking Core API (Ledger Endpoints) at $API_URL"
echo "=================================================="

# Step 1: Authentication
echo "📝 Step 1: Authenticating with admin credentials"
echo "--------------------------------------------------"

LOGIN_RESULT=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password"
  }')

# Extract token
TOKEN=$(echo "$LOGIN_RESULT" | grep -o '"token":"[^"]*' | grep -o '[^"]*$')

if [ -z "$TOKEN" ]; then
  echo "❌ Authentication failed"
  pretty_json "$LOGIN_RESULT"
  exit 1
else
  echo "✅ Authentication successful"
  echo "Token: ${TOKEN}"
fi

echo

# Step 2: Create Ledger Account #1
echo "📝 Step 2: Creating the first ledger account (balance=1000)"
echo "--------------------------------------------------"

CREATE_ACC1_RESULT=$(curl -s -X POST "$API_URL/ledger/account?balance=1000" \
  -H "Authorization: Bearer $TOKEN")

if [[ $CREATE_ACC1_RESULT == *"account_id"* ]]; then
  echo "✅ First account created successfully"
  pretty_json "$CREATE_ACC1_RESULT"
  # Extract the account ID
  ACC1_ID=$(echo "$CREATE_ACC1_RESULT" | grep -o '"account_id":"[^"]*' | grep -o '[^"]*$')
  echo "Account 1 ID: $ACC1_ID"
else
  echo "❌ Failed to create the first ledger account"
  pretty_json "$CREATE_ACC1_RESULT"
  exit 1
fi

echo

# Step 3: Create Ledger Account #2
echo "📝 Step 3: Creating the second ledger account (balance=500)"
echo "--------------------------------------------------"

CREATE_ACC2_RESULT=$(curl -s -X POST "$API_URL/ledger/account?balance=500" \
  -H "Authorization: Bearer $TOKEN")

if [[ $CREATE_ACC2_RESULT == *"account_id"* ]]; then
  echo "✅ Second account created successfully"
  pretty_json "$CREATE_ACC2_RESULT"
  # Extract the account ID
  ACC2_ID=$(echo "$CREATE_ACC2_RESULT" | grep -o '"account_id":"[^"]*' | grep -o '[^"]*$')
  echo "Account 2 ID: $ACC2_ID"
else
  echo "❌ Failed to create the second ledger account"
  pretty_json "$CREATE_ACC2_RESULT"
  exit 1
fi

echo

# Step 4: Transfer Funds
echo "📝 Step 4: Transferring 300 from the first account to the second"
echo "--------------------------------------------------"

TRANSFER_RESULT=$(curl -s -X POST "$API_URL/ledger/transfer?from=$ACC1_ID&to=$ACC2_ID&amount=300" \
  -H "Authorization: Bearer $TOKEN")

if [[ $TRANSFER_RESULT == *"Transfer successful"* ]]; then
  echo "✅ Transfer completed successfully"
  pretty_json "$TRANSFER_RESULT"
else
  echo "❌ Failed to transfer funds"
  pretty_json "$TRANSFER_RESULT"
  exit 1
fi

echo
echo "🎉 All Ledger endpoint tests completed successfully!"