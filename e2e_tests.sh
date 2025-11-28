#!/bin/bash

# E2E Testing Script for Animal Foundation CRM Backend

BASE_URL="http://localhost:8080"
API_PREFIX="/api/v1"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Global variables for storing IDs
ACCESS_TOKEN=""
REFRESH_TOKEN=""
ANIMAL_ID=""

# Test result tracking
declare -a TEST_RESULTS

# Function to print colored output
print_header() {
    echo -e "\n${BLUE}============================================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}============================================================${NC}\n"
}

print_test() {
    local name="$1"
    local status="$2"
    local details="$3"

    ((TOTAL_TESTS++))

    if [ "$status" == "PASS" ]; then
        echo -e "${GREEN}✓ $name${NC} $details"
        ((PASSED_TESTS++))
        TEST_RESULTS+=("PASS|$name")
    elif [ "$status" == "FAIL" ]; then
        echo -e "${RED}✗ $name${NC} $details"
        ((FAILED_TESTS++))
        TEST_RESULTS+=("FAIL|$name")
    else
        echo -e "${YELLOW}⊘ $name${NC} $details"
        TEST_RESULTS+=("SKIP|$name")
    fi
}

# Function to make HTTP requests using curl
make_request() {
    local method="$1"
    local endpoint="$2"
    local data="$3"
    local token="$4"
    local url="${BASE_URL}${API_PREFIX}${endpoint}"

    # If no token is provided, use the global ACCESS_TOKEN
    if [ -z "$token" ]; then
        token="$ACCESS_TOKEN"
    fi

    local response
    if [ -n "$data" ]; then
        response=$(curl -s -w "%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $token" \
            -d "$data" \
            "$url")
    else
        response=$(curl -s -w "%{http_code}" -X "$method" \
            -H "Authorization: Bearer $token" \
            "$url")
    fi

    http_code="${response: -3}"
    body="${response:0:${#response}-3}"

    if [[ "$http_code" -ge 200 && "$http_code" -lt 300 ]]; then
        echo "SUCCESS"
        echo "$body"
    else
        echo "FAILED"
        echo "$http_code"
        echo "$body"
    fi
}

# ============================================================================
# AUTHENTICATION TESTS
# ============================================================================

test_authentication() {
    print_header "TESTING: Authentication"

    # Test 1: Login
    echo "Logging in..."
    response=$(make_request "POST" "/auth/login" '{
        "email": "sarah.johnson@happypaws.org",
        "password": "password123"
    }')

    if echo "$response" | grep -q "SUCCESS"; then
        ACCESS_TOKEN=$(echo "$response" | jq -r .access_token)
        REFRESH_TOKEN=$(echo "$response" | jq -r .refresh_token)
        if [ -n "$ACCESS_TOKEN" ] && [ -n "$REFRESH_TOKEN" ]; then
            print_test "POST /auth/login" "PASS" "- Token obtained"
        else
            print_test "POST /auth/login" "FAIL" "- No token in response"
            return 1
        fi
    else
        print_test "POST /auth/login" "FAIL" "- Login failed"
        return 1
    fi
}

# ============================================================================
# ANIMAL MANAGEMENT TESTS
# ============================================================================

test_animals() {
    print_header "TESTING: Animal Management"

    # Test 1: Create animal
    response=$(make_request "POST" "/animals" '{
        "name": {
            "en": "Max",
            "pl": "Maks"
        },
        "category": "mammal",
        "species": "dog",
        "breed": "Golden Retriever",
        "sex": "male",
        "status": "available",
        "date_of_birth": "2022-05-15T00:00:00Z",
        "age_estimated": false,
        "color": "golden",
        "size": "large",
        "weight": 30.5,
        "description": {
            "en": "Max is a friendly Golden Retriever..."
        },
        "medical": {
            "vaccinated": true,
            "sterilized": true,
            "microchipped": true,
            "microchip_number": "123456789012345",
            "health_status": "healthy"
        },
        "behavior": {
            "temperament": ["friendly", "playful"],
            "good_with_kids": true,
            "good_with_dogs": true,
            "good_with_cats": false,
            "house_trained": true
        },
        "intake_date": "2024-01-15T00:00:00Z",
        "intake_reason": "Owner surrender",
        "location": "Kennel 5",
        "adoption_fee": 150.00
    }')

    if echo "$response" | grep -q "SUCCESS"; then
        ANIMAL_ID=$(echo "$response" | jq -r .id)
        print_test "POST /animals" "PASS" "- ID: $ANIMAL_ID"
    else
        print_test "POST /animals" "FAIL"
        return 1
    fi

    # Test 2: Get animal by ID
    if [ -n "$ANIMAL_ID" ]; then
        response=$(make_request "GET" "/animals/$ANIMAL_ID" "")
        if echo "$response" | grep -q "SUCCESS"; then
            print_test "GET /animals/:id" "PASS"
        else
            print_test "GET /animals/:id" "FAIL"
        fi
    fi

    # Test 3: List animals
    response=$(make_request "GET" "/animals" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /animals" "PASS"
    else
        print_test "GET /animals" "FAIL"
    fi

    # Test 4: Search animals
    response=$(make_request "GET" "/animals?search=Max" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /animals?search=Max" "PASS"
    else
        print_test "GET /animals?search=Max" "FAIL"
    fi

    # Test 5: Paginate animals
    response=$(make_request "GET" "/animals?limit=5&offset=5" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /animals?limit=5&offset=5" "PASS"
    else
        print_test "GET /animals?limit=5&offset=5" "FAIL"
    fi

    # Test 6: Update animal
    if [ -n "$ANIMAL_ID" ]; then
        response=$(make_request "PUT" "/animals/$ANIMAL_ID" '{
            "weight": 31.0
        }')
        if echo "$response" | grep -q "SUCCESS"; then
            print_test "PUT /animals/:id" "PASS"
        else
            print_test "PUT /animals/:id" "FAIL"
        fi
    fi

    # Test 7: Delete animal
    if [ -n "$ANIMAL_ID" ]; then
        response=$(make_request "DELETE" "/animals/$ANIMAL_ID" "")
        if echo "$response" | grep -q "SUCCESS"; then
            print_test "DELETE /animals/:id" "PASS"
        else
            print_test "DELETE /animals/:id" "FAIL"
        fi
    fi
}

# ============================================================================
# SESSION MANAGEMENT TESTS
# ============================================================================

test_session_management() {
    print_header "TESTING: Session Management"

    # Test 1: API behavior with an expired token
    echo "Generating expired token..."
    response=$(make_request "POST" "/test/generate-token" "{\"expires_in\": 1}")
    if echo "$response" | grep -q "SUCCESS"; then
        expired_token=$(echo "$response" | jq -r .access_token)
        sleep 2
        response=$(make_request "GET" "/animals" "" "$expired_token")
        if echo "$response" | grep -q "FAILED"; then
            print_test "GET /animals with expired token" "PASS"
        else
            print_test "GET /animals with expired token" "FAIL"
        fi
    else
        print_test "POST /test/generate-token" "FAIL"
    fi

    # Test 2: The token refresh mechanism
    response=$(make_request "POST" "/auth/refresh" "{\"refresh_token\": \"$REFRESH_TOKEN\"}")
    if echo "$response" | grep -q "SUCCESS"; then
        ACCESS_TOKEN=$(echo "$response" | jq -r .access_token)
        if [ -n "$ACCESS_TOKEN" ]; then
            print_test "POST /auth/refresh" "PASS" "- New token obtained"
        else
            print_test "POST /auth/refresh" "FAIL" "- No new token in response"
        fi
    else
        print_test "POST /auth/refresh" "FAIL" "- Token refresh failed"
    fi

    # Test 3: Rejection of stale/invalid tokens
    response=$(make_request "GET" "/animals" "" "invalid_token")
    if echo "$response" | grep -q "FAILED"; then
        print_test "GET /animals with invalid token" "PASS"
    else
        print_test "GET /animals with invalid token" "FAIL"
    fi
}

# ============================================================================
# MAIN EXECUTION
# ============================================================================

main() {
    echo -e "\n${BLUE}╔═══════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║       Animal Foundation CRM - E2E API Testing             ║${NC}"
    echo -e "${BLUE}╚═══════════════════════════════════════════════════════════╝${NC}"

    echo -e "\nTesting against: ${BLUE}${BASE_URL}${NC}\n"

    # Run all test modules
    test_authentication

    # Only continue if authentication succeeded
    if [ -n "$ACCESS_TOKEN" ]; then
        test_animals
        test_session_management
    else
        echo -e "${RED}Authentication failed. Cannot continue with other tests.${NC}"
    fi

    # Exit with appropriate code
    if [ $FAILED_TESTS -gt 0 ]; then
        exit 1
    else
        exit 0
    fi
}

# Run main function
main
