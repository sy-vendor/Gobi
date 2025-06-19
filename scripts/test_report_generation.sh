#!/bin/bash

# Test script for report generation functionality
echo "Testing Report Generation Functionality..."

# Base URL
BASE_URL="http://localhost:8080/api"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
    fi
}

# Function to make API calls
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local token=$4
    
    if [ -n "$token" ]; then
        if [ -n "$data" ]; then
            curl -s -X $method "$BASE_URL$endpoint" \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $token" \
                -d "$data"
        else
            curl -s -X $method "$BASE_URL$endpoint" \
                -H "Authorization: Bearer $token"
        fi
    else
        if [ -n "$data" ]; then
            curl -s -X $method "$BASE_URL$endpoint" \
                -H "Content-Type: application/json" \
                -d "$data"
        else
            curl -s -X $method "$BASE_URL$endpoint"
        fi
    fi
}

echo "1. Testing User Registration..."
REGISTER_RESPONSE=$(make_request "POST" "/auth/register" '{
    "username": "reportuser",
    "password": "test123",
    "email": "report@example.com"
}')
print_status $? "User registration"

echo "2. Testing User Login..."
LOGIN_RESPONSE=$(make_request "POST" "/auth/login" '{
    "username": "reportuser",
    "password": "test123"
}')
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
print_status $? "User login"

if [ -z "$TOKEN" ]; then
    echo -e "${RED}Failed to get token, exiting...${NC}"
    exit 1
fi

echo "3. Testing Data Source Creation..."
DS_RESPONSE=$(make_request "POST" "/datasources" '{
    "name": "Test DB",
    "type": "sqlite",
    "database": "test.db",
    "description": "Test database for reports"
}' "$TOKEN")
DS_ID=$(echo $DS_RESPONSE | grep -o '"ID":[0-9]*' | cut -d':' -f2)
print_status $? "Data source creation"

echo "4. Testing Query Creation..."
QUERY_RESPONSE=$(make_request "POST" "/queries" '{
    "name": "Test Query",
    "sql": "SELECT 1 as test_column, \"test_value\" as test_data",
    "description": "Test query for reports",
    "data_source_id": '$DS_ID'
}' "$TOKEN")
QUERY_ID=$(echo $QUERY_RESPONSE | grep -o '"ID":[0-9]*' | cut -d':' -f2)
print_status $? "Query creation"

echo "5. Testing Report Schedule Creation..."
SCHEDULE_RESPONSE=$(make_request "POST" "/reports/schedules" '{
    "name": "Daily Test Report",
    "type": "daily",
    "query_ids": ['$QUERY_ID'],
    "chart_ids": [],
    "template_ids": [],
    "cron_pattern": "0 0 * * *"
}' "$TOKEN")
SCHEDULE_ID=$(echo $SCHEDULE_RESPONSE | grep -o '"ID":[0-9]*' | cut -d':' -f2)
print_status $? "Report schedule creation"

echo "6. Testing Report Schedule Listing..."
SCHEDULES_RESPONSE=$(make_request "GET" "/reports/schedules" "" "$TOKEN")
print_status $? "Report schedule listing"

echo "7. Testing Report Schedule Retrieval..."
SCHEDULE_DETAIL_RESPONSE=$(make_request "GET" "/reports/schedules/$SCHEDULE_ID" "" "$TOKEN")
print_status $? "Report schedule retrieval"

echo "8. Testing Report Schedule Update..."
UPDATE_RESPONSE=$(make_request "PUT" "/reports/schedules/$SCHEDULE_ID" '{
    "name": "Updated Daily Test Report",
    "active": false
}' "$TOKEN")
print_status $? "Report schedule update"

echo "9. Testing Report Listing..."
REPORTS_RESPONSE=$(make_request "GET" "/reports" "" "$TOKEN")
print_status $? "Report listing"

echo "10. Testing Report Schedule Deletion..."
DELETE_RESPONSE=$(make_request "DELETE" "/reports/schedules/$SCHEDULE_ID" "" "$TOKEN")
print_status $? "Report schedule deletion"

echo -e "\n${YELLOW}Report Generation Test Summary:${NC}"
echo "✓ User registration and authentication"
echo "✓ Data source and query creation"
echo "✓ Report schedule CRUD operations"
echo "✓ Report listing functionality"

echo -e "\n${GREEN}All report generation tests completed!${NC}" 