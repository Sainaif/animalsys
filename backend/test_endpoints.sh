#!/bin/bash

# Comprehensive API Testing Script for Animal Foundation CRM Backend
# Tests all critical endpoints with sample data

BASE_URL="http://localhost:8081"
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
ANIMAL_ID=""
CONDITION_ID=""
MEDICATION_ID=""
TREATMENT_PLAN_ID=""
ADOPTION_ID=""
DONOR_ID=""
DONATION_ID=""
CAMPAIGN_ID=""
EVENT_ID=""
VOLUNTEER_ID=""
INVENTORY_ID=""

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

# Function to make HTTP requests using wget
make_request() {
    local method="$1"
    local endpoint="$2"
    local data="$3"
    local url="${BASE_URL}${API_PREFIX}${endpoint}"
    local temp_file="/tmp/api_response_$$.json"
    local headers_file="/tmp/api_headers_$$.txt"

    # Build wget command
    local wget_cmd="wget -q -O $temp_file -S"

    # Add authorization header if token exists
    if [ -n "$ACCESS_TOKEN" ]; then
        wget_cmd="$wget_cmd --header='Authorization: Bearer $ACCESS_TOKEN'"
    fi

    # Add content type for POST/PUT
    if [ "$method" == "POST" ] || [ "$method" == "PUT" ]; then
        wget_cmd="$wget_cmd --header='Content-Type: application/json'"
    fi

    # Set method and data
    if [ "$method" == "POST" ]; then
        wget_cmd="$wget_cmd --post-data='$data'"
    elif [ "$method" == "PUT" ]; then
        wget_cmd="$wget_cmd --method=PUT --body-data='$data'"
    elif [ "$method" == "DELETE" ]; then
        wget_cmd="$wget_cmd --method=DELETE"
    fi

    wget_cmd="$wget_cmd '$url' 2>&1"

    # Execute request
    eval $wget_cmd > $headers_file

    # Check if request was successful
    if grep -q "HTTP/1.1 2[0-9][0-9]" $headers_file 2>/dev/null; then
        echo "SUCCESS"
        if [ -f "$temp_file" ]; then
            cat "$temp_file"
        fi
    else
        echo "FAILED"
        cat $headers_file 2>/dev/null | grep "HTTP/1.1" | head -1
    fi

    # Cleanup
    rm -f "$temp_file" "$headers_file" 2>/dev/null
}

# ============================================================================
# AUTHENTICATION TESTS
# ============================================================================

test_authentication() {
    print_header "TESTING: Authentication & User Management"

    # Test 1: Register (might fail if user exists, that's ok)
    echo "Attempting to register admin user..."
    response=$(make_request "POST" "/auth/register" '{
        "email": "admin@animalsys.com",
        "password": "AdminPass123!",
        "first_name": "System",
        "last_name": "Administrator",
        "role": "super_admin"
    }')

    if echo "$response" | grep -q "SUCCESS"; then
        print_test "POST /auth/register" "PASS" "- User registered"
    else
        print_test "POST /auth/register" "SKIP" "- User may already exist"
    fi

    # Test 2: Login (critical)
    echo "Logging in..."
    response=$(make_request "POST" "/auth/login" '{
        "email": "admin@animalsys.com",
        "password": "AdminPass123!"
    }')

    if echo "$response" | grep -q "SUCCESS"; then
        # Extract access token
        ACCESS_TOKEN=$(echo "$response" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
        if [ -n "$ACCESS_TOKEN" ]; then
            print_test "POST /auth/login" "PASS" "- Token obtained"
        else
            print_test "POST /auth/login" "FAIL" "- No token in response"
            return 1
        fi
    else
        print_test "POST /auth/login" "FAIL" "- Login failed"
        return 1
    fi

    # Test 3: Get current user
    response=$(make_request "GET" "/auth/me" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /auth/me" "PASS"
    else
        print_test "GET /auth/me" "FAIL"
    fi

    # Test 4: List users
    response=$(make_request "GET" "/users?limit=10" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /users" "PASS"
    else
        print_test "GET /users" "FAIL"
    fi
}

# ============================================================================
# ANIMAL MANAGEMENT TESTS
# ============================================================================

test_animals() {
    print_header "TESTING: Animal Management"

    # Test 1: Create animal
    response=$(make_request "POST" "/animals" '{
        "name": {"en": "Max", "pl": "Maks"},
        "species": "dog",
        "breed": "Golden Retriever",
        "category": "mammal",
        "age_years": 3,
        "age_months": 6,
        "gender": "male",
        "color": "Golden",
        "weight_kg": 28.5,
        "status": "available",
        "location": "Shelter A",
        "description": {
            "en": "Friendly and energetic dog",
            "pl": "Przyjazny i energiczny pies"
        }
    }')

    if echo "$response" | grep -q "SUCCESS"; then
        ANIMAL_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
        print_test "POST /animals" "PASS" "- ID: ${ANIMAL_ID:0:8}..."
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
    response=$(make_request "GET" "/animals?limit=10" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /animals" "PASS"
    else
        print_test "GET /animals" "FAIL"
    fi

    # Test 4: Get available animals
    response=$(make_request "GET" "/animals/available" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /animals/available" "PASS"
    else
        print_test "GET /animals/available" "FAIL"
    fi

    # Test 5: Get animals by species
    response=$(make_request "GET" "/animals/species/dog" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /animals/species/:species" "PASS"
    else
        print_test "GET /animals/species/:species" "FAIL"
    fi

    # Test 6: Update animal
    if [ -n "$ANIMAL_ID" ]; then
        response=$(make_request "PUT" "/animals/$ANIMAL_ID" '{
            "weight_kg": 29.0,
            "location": "Shelter B"
        }')
        if echo "$response" | grep -q "SUCCESS"; then
            print_test "PUT /animals/:id" "PASS"
        else
            print_test "PUT /animals/:id" "FAIL"
        fi
    fi
}

# ============================================================================
# MEDICAL RECORDS TESTS
# ============================================================================

test_medical() {
    print_header "TESTING: Medical Records & Treatment"

    if [ -z "$ANIMAL_ID" ]; then
        print_test "Medical Tests" "SKIP" "- No animal ID available"
        return
    fi

    # Test 1: Create medical condition
    response=$(make_request "POST" "/medical-conditions" "{
        \"animal_id\": \"$ANIMAL_ID\",
        \"name\": \"Minor skin irritation\",
        \"diagnosis\": \"Allergic dermatitis\",
        \"severity\": \"mild\",
        \"status\": \"active\",
        \"symptoms\": [\"Redness\", \"Itching\"],
        \"diagnosis_date\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
        \"notes\": \"Monitor for improvement\"
    }")

    if echo "$response" | grep -q "SUCCESS"; then
        CONDITION_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
        print_test "POST /medical-conditions" "PASS" "- ID: ${CONDITION_ID:0:8}..."
    else
        print_test "POST /medical-conditions" "FAIL"
        return
    fi

    # Test 2: List medical conditions
    response=$(make_request "GET" "/medical-conditions?limit=10" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /medical-conditions" "PASS"
    else
        print_test "GET /medical-conditions" "FAIL"
    fi

    # Test 3: Create medication
    response=$(make_request "POST" "/medications" "{
        \"animal_id\": \"$ANIMAL_ID\",
        \"condition_id\": \"$CONDITION_ID\",
        \"name\": \"Antihistamine Cream\",
        \"dosage\": \"Apply twice daily\",
        \"frequency\": \"2x daily\",
        \"route\": \"topical\",
        \"start_date\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
        \"status\": \"active\",
        \"instructions\": \"Apply thin layer to affected areas\"
    }")

    if echo "$response" | grep -q "SUCCESS"; then
        MEDICATION_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
        print_test "POST /medications" "PASS" "- ID: ${MEDICATION_ID:0:8}..."
    else
        print_test "POST /medications" "FAIL"
    fi

    # Test 4: List medications
    response=$(make_request "GET" "/medications?limit=10" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /medications" "PASS"
    else
        print_test "GET /medications" "FAIL"
    fi

    # Test 5: Create treatment plan
    response=$(make_request "POST" "/treatment-plans" "{
        \"animal_id\": \"$ANIMAL_ID\",
        \"name\": \"Dermatitis Treatment Plan\",
        \"description\": \"7-day treatment for allergic dermatitis\",
        \"start_date\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
        \"status\": \"active\",
        \"goals\": [\"Reduce inflammation\", \"Eliminate itching\"]
    }")

    if echo "$response" | grep -q "SUCCESS"; then
        print_test "POST /treatment-plans" "PASS"
    else
        print_test "POST /treatment-plans" "FAIL"
    fi
}

# ============================================================================
# ADOPTION TESTS
# ============================================================================

test_adoptions() {
    print_header "TESTING: Adoption Management"

    if [ -z "$ANIMAL_ID" ]; then
        print_test "Adoption Tests" "SKIP" "- No animal ID available"
        return
    fi

    # Test 1: Create adoption application
    response=$(make_request "POST" "/adoptions" "{
        \"animal_id\": \"$ANIMAL_ID\",
        \"applicant_name\": \"John Doe\",
        \"applicant_email\": \"john.doe@example.com\",
        \"applicant_phone\": \"+1234567890\",
        \"address\": \"123 Main St, City, State 12345\",
        \"housing_type\": \"house\",
        \"has_yard\": true,
        \"has_other_pets\": false,
        \"experience_level\": \"intermediate\",
        \"status\": \"pending\",
        \"application_date\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
    }")

    if echo "$response" | grep -q "SUCCESS"; then
        ADOPTION_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
        print_test "POST /adoptions" "PASS" "- ID: ${ADOPTION_ID:0:8}..."
    else
        print_test "POST /adoptions" "FAIL"
        return
    fi

    # Test 2: List adoptions
    response=$(make_request "GET" "/adoptions?limit=10" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /adoptions" "PASS"
    else
        print_test "GET /adoptions" "FAIL"
    fi

    # Test 3: Get pending adoptions
    response=$(make_request "GET" "/adoptions/pending" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /adoptions/pending" "PASS"
    else
        print_test "GET /adoptions/pending" "FAIL"
    fi
}

# ============================================================================
# DONATION TESTS
# ============================================================================

test_donations() {
    print_header "TESTING: Donations & Fundraising"

    # Test 1: Create donor
    response=$(make_request "POST" "/donors" '{
        "first_name": "Jane",
        "last_name": "Smith",
        "email": "jane.smith@example.com",
        "phone": "+1987654321",
        "address": "456 Oak Ave, City, State 12345",
        "donor_type": "individual",
        "status": "active"
    }')

    if echo "$response" | grep -q "SUCCESS"; then
        DONOR_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
        print_test "POST /donors" "PASS" "- ID: ${DONOR_ID:0:8}..."
    else
        print_test "POST /donors" "FAIL"
        return
    fi

    # Test 2: Create donation
    response=$(make_request "POST" "/donations" "{
        \"donor_id\": \"$DONOR_ID\",
        \"amount\": 100.00,
        \"currency\": \"USD\",
        \"donation_type\": \"monetary\",
        \"payment_method\": \"credit_card\",
        \"status\": \"completed\",
        \"donation_date\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
    }")

    if echo "$response" | grep -q "SUCCESS"; then
        print_test "POST /donations" "PASS"
    else
        print_test "POST /donations" "FAIL"
    fi

    # Test 3: List donations
    response=$(make_request "GET" "/donations?limit=10" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /donations" "PASS"
    else
        print_test "GET /donations" "FAIL"
    fi

    # Test 4: Create campaign
    response=$(make_request "POST" "/campaigns" '{
        "name": "Winter Care Campaign",
        "description": "Raising funds for winter animal care",
        "goal_amount": 5000.00,
        "currency": "USD",
        "status": "active"
    }')

    if echo "$response" | grep -q "SUCCESS"; then
        print_test "POST /campaigns" "PASS"
    else
        print_test "POST /campaigns" "FAIL"
    fi

    # Test 5: List campaigns
    response=$(make_request "GET" "/campaigns?limit=10" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /campaigns" "PASS"
    else
        print_test "GET /campaigns" "FAIL"
    fi
}

# ============================================================================
# EVENTS & VOLUNTEERS TESTS
# ============================================================================

test_events_volunteers() {
    print_header "TESTING: Events & Volunteers"

    # Test 1: Create event
    response=$(make_request "POST" "/events" '{
        "name": "Adoption Fair",
        "description": "Meet adoptable animals",
        "event_type": "adoption_event",
        "location": "City Park",
        "capacity": 100,
        "status": "scheduled"
    }')

    if echo "$response" | grep -q "SUCCESS"; then
        print_test "POST /events" "PASS"
    else
        print_test "POST /events" "FAIL"
    fi

    # Test 2: List events
    response=$(make_request "GET" "/events?limit=10" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /events" "PASS"
    else
        print_test "GET /events" "FAIL"
    fi

    # Test 3: Create volunteer
    response=$(make_request "POST" "/volunteers" '{
        "first_name": "Alice",
        "last_name": "Johnson",
        "email": "alice.johnson@example.com",
        "phone": "+1122334455",
        "skills": ["Animal care", "Event organization"],
        "availability": "weekends",
        "status": "active"
    }')

    if echo "$response" | grep -q "SUCCESS"; then
        print_test "POST /volunteers" "PASS"
    else
        print_test "POST /volunteers" "FAIL"
    fi

    # Test 4: List volunteers
    response=$(make_request "GET" "/volunteers?limit=10" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /volunteers" "PASS"
    else
        print_test "GET /volunteers" "FAIL"
    fi
}

# ============================================================================
# INVENTORY TESTS
# ============================================================================

test_inventory() {
    print_header "TESTING: Inventory Management"

    # Test 1: Create inventory item
    response=$(make_request "POST" "/inventory" '{
        "name": "Dog Food - Premium",
        "sku": "DF-PREM-001",
        "category": "food",
        "quantity": 50,
        "unit": "bags",
        "unit_cost": 25.00,
        "reorder_level": 10,
        "status": "active"
    }')

    if echo "$response" | grep -q "SUCCESS"; then
        INVENTORY_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
        print_test "POST /inventory" "PASS" "- ID: ${INVENTORY_ID:0:8}..."
    else
        print_test "POST /inventory" "FAIL"
        return
    fi

    # Test 2: List inventory
    response=$(make_request "GET" "/inventory?limit=10" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /inventory" "PASS"
    else
        print_test "GET /inventory" "FAIL"
    fi

    # Test 3: Add stock
    if [ -n "$INVENTORY_ID" ]; then
        response=$(make_request "POST" "/inventory/$INVENTORY_ID/add-stock" '{
            "quantity": 20,
            "notes": "Monthly restock"
        }')
        if echo "$response" | grep -q "SUCCESS"; then
            print_test "POST /inventory/:id/add-stock" "PASS"
        else
            print_test "POST /inventory/:id/add-stock" "FAIL"
        fi
    fi
}

# ============================================================================
# MONITORING TESTS
# ============================================================================

test_monitoring() {
    print_header "TESTING: Monitoring & Audit"

    # Test 1: Get system health
    response=$(make_request "GET" "/monitoring/health" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /monitoring/health" "PASS"
    else
        print_test "GET /monitoring/health" "FAIL"
    fi

    # Test 2: Get statistics
    response=$(make_request "GET" "/monitoring/statistics" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /monitoring/statistics" "PASS"
    else
        print_test "GET /monitoring/statistics" "FAIL"
    fi

    # Test 3: List audit logs
    response=$(make_request "GET" "/audit-logs?limit=10" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /audit-logs" "PASS"
    else
        print_test "GET /audit-logs" "FAIL"
    fi
}

# ============================================================================
# PUBLIC ENDPOINTS TESTS
# ============================================================================

test_public() {
    print_header "TESTING: Public Endpoints"

    # Temporarily clear access token for public endpoints
    local saved_token="$ACCESS_TOKEN"
    ACCESS_TOKEN=""

    # Test 1: Get species list
    response=$(make_request "GET" "/public/animals/species" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /public/animals/species" "PASS"
    else
        print_test "GET /public/animals/species" "FAIL"
    fi

    # Test 2: List public animals
    response=$(make_request "GET" "/public/animals?limit=10" "")
    if echo "$response" | grep -q "SUCCESS"; then
        print_test "GET /public/animals" "PASS"
    else
        print_test "GET /public/animals" "FAIL"
    fi

    # Restore access token
    ACCESS_TOKEN="$saved_token"
}

# ============================================================================
# GENERATE REPORT
# ============================================================================

generate_report() {
    print_header "TEST SUMMARY REPORT"

    echo "Total Tests: $TOTAL_TESTS"
    echo -e "${GREEN}Passed: $PASSED_TESTS${NC}"
    echo -e "${RED}Failed: $FAILED_TESTS${NC}"

    if [ $TOTAL_TESTS -gt 0 ]; then
        success_rate=$(( PASSED_TESTS * 100 / TOTAL_TESTS ))
        echo "Success Rate: ${success_rate}%"
    fi

    # Save report to file
    {
        echo "Animal Foundation CRM - API Test Report"
        echo "========================================"
        echo "Date: $(date)"
        echo ""
        echo "Summary:"
        echo "  Total Tests: $TOTAL_TESTS"
        echo "  Passed: $PASSED_TESTS"
        echo "  Failed: $FAILED_TESTS"
        echo "  Success Rate: ${success_rate}%"
        echo ""
        echo "Detailed Results:"
        for result in "${TEST_RESULTS[@]}"; do
            status="${result%%|*}"
            name="${result##*|}"
            echo "  [$status] $name"
        done
    } > test_report.txt

    echo -e "\n${BLUE}Detailed report saved to: test_report.txt${NC}"

    if [ $FAILED_TESTS -gt 0 ]; then
        echo -e "\n${RED}Some tests failed. Review the output above for details.${NC}"
    else
        echo -e "\n${GREEN}All tests passed!${NC}"
    fi
}

# ============================================================================
# MAIN EXECUTION
# ============================================================================

main() {
    echo -e "\n${BLUE}╔═══════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║   Animal Foundation CRM - Comprehensive API Testing      ║${NC}"
    echo -e "${BLUE}║              All Critical Endpoints Test Suite            ║${NC}"
    echo -e "${BLUE}╚═══════════════════════════════════════════════════════════╝${NC}"

    echo -e "\nTesting against: ${BLUE}${BASE_URL}${NC}\n"

    # Run all test modules
    test_authentication

    # Only continue if authentication succeeded
    if [ -n "$ACCESS_TOKEN" ]; then
        test_animals
        test_medical
        test_adoptions
        test_donations
        test_events_volunteers
        test_inventory
        test_monitoring
        test_public
    else
        echo -e "${RED}Authentication failed. Cannot continue with other tests.${NC}"
    fi

    # Generate report
    generate_report

    # Exit with appropriate code
    if [ $FAILED_TESTS -gt 0 ]; then
        exit 1
    else
        exit 0
    fi
}

# Run main function
main
