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
echo "🚀 Testing Banking Core API (Business Endpoints with KYC) at $API_URL"
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

# Step 2: Create Business Entity with additional KYC fields
echo "📝 Step 2: Creating a new business entity with KYC details"
echo "--------------------------------------------------"

CREATE_RESULT=$(curl -s -X POST "$API_URL/entities/business" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Acme Corporation",
    "registration_number": "ACME-123456",
    "address": "456 Corporate Blvd",
    "country": "US",
    "tax_id": "TAX-7890",
    "kyc_document_url": "http://example.com/business-doc.pdf"
  }')

# Extract business ID
BUSINESS_ID=$(echo "$CREATE_RESULT" | grep -o '"id":"[^"]*' | grep -o '[^"]*$')

if [ -z "$BUSINESS_ID" ]; then
  echo "❌ Failed to create business entity"
  pretty_json "$CREATE_RESULT"
  exit 1
else
  echo "✅ Business entity created successfully"
  echo "Business ID: $BUSINESS_ID"
  pretty_json "$CREATE_RESULT"
fi

echo

# Step 3: Retrieve Business Entity
echo "📝 Step 3: Retrieving the created business entity"
echo "--------------------------------------------------"

GET_RESULT=$(curl -s -X GET "$API_URL/entities/business/$BUSINESS_ID" \
  -H "Authorization: Bearer $TOKEN")

if [[ $GET_RESULT == *"name"* ]]; then
  echo "✅ Business entity retrieved successfully"
  pretty_json "$GET_RESULT"
else
  echo "❌ Failed to retrieve business entity"
  pretty_json "$GET_RESULT"
  exit 1
fi

echo

# Step 4: Update Business Entity (update KYC status to verified)
echo "📝 Step 4: Updating the business entity with verified KYC status"
echo "--------------------------------------------------"

UPDATE_RESULT=$(curl -s -X PATCH "$API_URL/entities/business/$BUSINESS_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "kyc_status": "verified"
  }')

if [[ $UPDATE_RESULT == *"verified"* ]]; then
  echo "✅ Business entity updated successfully"
  pretty_json "$UPDATE_RESULT"
else
  echo "❌ Failed to update business entity"
  pretty_json "$UPDATE_RESULT"
  exit 1
fi

echo

# Step 5: List Business Entities
echo "📝 Step 5: Listing all business entities"
echo "--------------------------------------------------"

LIST_RESULT=$(curl -s -X GET "$API_URL/entities/business?limit=10" \
  -H "Authorization: Bearer $TOKEN")

if [[ $LIST_RESULT == \[* ]]; then
  COUNT=$(echo "$LIST_RESULT" | grep -o '"id"' | wc -l)
  echo "✅ Listed $COUNT business entities successfully"
  if [ $COUNT -gt 1 ]; then
    FIRST_ENTITY=$(echo "$LIST_RESULT" | sed 's/\[//' | sed 's/,{.*$//')
    pretty_json "$FIRST_ENTITY"
    echo "... and $(($COUNT - 1)) more business entities"
  else
    pretty_json "$LIST_RESULT"
  fi
else
  echo "❌ Failed to list business entities"
  pretty_json "$LIST_RESULT"
  exit 1
fi

echo
echo "🎉 All Business endpoint tests completed successfully!"