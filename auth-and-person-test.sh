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
echo "🚀 Testing Banking Core API (Person KYC) at $API_URL"
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

# Step 2: Create Person Entity with additional KYC fields
echo "📝 Step 2: Creating a new person entity with KYC details"
echo "--------------------------------------------------"

CREATE_RESULT=$(curl -s -X POST "$API_URL/entities/person" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "date_of_birth": "1990-01-15",
    "email": "john.doe@example.com",
    "phone_number": "+1-555-123-4567",
    "street1": "123 Main St",
    "city": "New York",
    "state": "NY",
    "postal_code": "10001",
    "country": "US",
    "government_id": "ABC123456",
    "nationality": "US",
    "kyc_document_url": "http://example.com/doc.pdf"
  }')

# Extract person ID
PERSON_ID=$(echo "$CREATE_RESULT" | grep -o '"id":"[^"]*' | grep -o '[^"]*$')

if [ -z "$PERSON_ID" ]; then
  echo "❌ Failed to create person entity"
  pretty_json "$CREATE_RESULT"
  exit 1
else
  echo "✅ Person entity created successfully"
  echo "Person ID: $PERSON_ID"
  pretty_json "$CREATE_RESULT"
fi

echo

# Step 3: Retrieve the created person entity and verify KYC fields
echo "📝 Step 3: Retrieving the created person entity"
echo "--------------------------------------------------"

GET_RESULT=$(curl -s -X GET "$API_URL/entities/person/$PERSON_ID" \
  -H "Authorization: Bearer $TOKEN")

if [[ $GET_RESULT == *"first_name"* ]]; then
  echo "✅ Person entity retrieved successfully"
  pretty_json "$GET_RESULT"
else
  echo "❌ Failed to retrieve person entity"
  pretty_json "$GET_RESULT"
  exit 1
fi

echo

# Step 4: Update Person Entity's KYC status to verified
echo "📝 Step 4: Updating the person entity (KYC status to verified)"
echo "--------------------------------------------------"

UPDATE_RESULT=$(curl -s -X PATCH "$API_URL/entities/person/$PERSON_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "kyc_status": "verified"
  }')

if [[ $UPDATE_RESULT == *"verified"* ]]; then
  echo "✅ Person entity updated successfully"
  pretty_json "$UPDATE_RESULT"
else
  echo "❌ Failed to update person entity"
  pretty_json "$UPDATE_RESULT"
  exit 1
fi

echo

# Step 5: List all person entities
echo "📝 Step 5: Listing all person entities"
echo "--------------------------------------------------"

LIST_RESULT=$(curl -s -X GET "$API_URL/entities/person?limit=10" \
  -H "Authorization: Bearer $TOKEN")

if [[ $LIST_RESULT == \[* ]]; then
  PEOPLE_COUNT=$(echo "$LIST_RESULT" | grep -o '"id"' | wc -l)
  echo "✅ Listed $PEOPLE_COUNT person entities successfully"
  
  # Only show the first item if there are multiple
  if [ $PEOPLE_COUNT -gt 1 ]; then
    FIRST_PERSON=$(echo "$LIST_RESULT" | sed 's/\[//' | sed 's/,{.*$//')
    pretty_json "$FIRST_PERSON"
    echo "... and $(($PEOPLE_COUNT-1)) more person entities"
  else
    pretty_json "$LIST_RESULT"
  fi
else
  echo "❌ Failed to list person entities"
  pretty_json "$LIST_RESULT"
  exit 1
fi

echo
echo "🎉 All Person KYC endpoint tests completed successfully!"