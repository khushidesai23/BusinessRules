#!/bin/bash

# Quick test setup and runner
# Usage: ./run_tests.sh [test-path]

set -e

echo "🧪 Business Rules Engine - Test Runner"
echo "======================================"

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker to run tests."
    exit 1
fi

# Parse arguments
TEST_PATH="${1:-.}"
TEST_NAME=$(basename "$TEST_PATH")

echo ""
echo "📦 Starting test database..."

# Start the test database container
if docker ps -a | grep -q "business-rules-test-db"; then
    echo "   (Removing existing container...)"
    docker stop business-rules-test-db 2>/dev/null || true
    docker rm business-rules-test-db 2>/dev/null || true
fi

docker-compose -f docker-compose.test.yml up -d

# Wait for database to be ready
echo "   (Waiting for database to be ready...)"
ready=0
for _ in $(seq 1 30); do
    if docker exec business-rules-test-db pg_isready -U postgres -d postgres >/dev/null 2>&1; then
        ready=1
        break
    fi
    sleep 1
done

if [ "$ready" -ne 1 ]; then
    echo "ERROR: Test database did not become ready in time."
    docker-compose -f docker-compose.test.yml down
    exit 1
fi

# Ensure the postgres password matches the DSN used by tests
configured=0
for _ in $(seq 1 30); do
    if docker exec business-rules-test-db psql -U postgres -d postgres -c "ALTER USER postgres PASSWORD 'postgres';" >/dev/null 2>&1; then
        configured=1
        break
    fi
    sleep 1
done

if [ "$configured" -ne 1 ]; then
    echo "ERROR: Failed to configure the test database credentials."
    docker-compose -f docker-compose.test.yml down
    exit 1
fi

# Set environment variable for tests
export POSTGRES_TEST_DSN="host=127.0.0.1 user=postgres password=postgres dbname=postgres port=55432 sslmode=disable"

echo "✅ Database ready"
echo ""

# Run tests
echo "🏃 Running tests..."
echo ""

if [ "$TEST_PATH" = "." ]; then
    echo "Running: go test ./... -v"
    go test ./... -v
else
    echo "Running: go test $TEST_PATH -v"
    go test "$TEST_PATH" -v
fi

TEST_RESULT=$?

echo ""
echo "======================================"

# Cleanup
echo "🧹 Cleaning up..."
docker-compose -f docker-compose.test.yml down

if [ $TEST_RESULT -eq 0 ]; then
    echo ""
    echo "✅ All tests passed!"
    exit 0
else
    echo ""
    echo "❌ Some tests failed"
    exit 1
fi
