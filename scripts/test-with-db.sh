#!/bin/bash

set -e

# Determine compose command
if command -v podman-compose &> /dev/null; then
    COMPOSE_CMD="podman-compose"
elif command -v docker-compose &> /dev/null; then
    COMPOSE_CMD="docker-compose"
else
    echo "Error: Neither podman-compose nor docker-compose found"
    exit 1
fi

# Clean up any existing containers
$COMPOSE_CMD -f docker-compose.test.yml down

# Start the test database with replacement
$COMPOSE_CMD -f docker-compose.test.yml up -d --force-recreate mariadb-test

# Wait for database to be ready (with timeout)
echo "Waiting for database to be ready..."
timeout=30
while ! $COMPOSE_CMD -f docker-compose.test.yml exec mariadb-test mysqladmin ping -h localhost --silent; do
    sleep 1
    timeout=$((timeout-1))
    if [ $timeout -le 0 ]; then
        echo "Timeout waiting for database"
        $COMPOSE_CMD -f docker-compose.test.yml logs mariadb-test
        exit 1
    fi
done

# Run tests with the test database
export MARIADB_TEST_DSN="testuser:testpass@tcp(localhost:3306)/testdb"
go test ./pkg/todo/... -v -coverprofile=coverage.out
test_exit=$?

# Generate coverage report if tests passed
if [ $test_exit -eq 0 ]; then
    go tool cover -func=coverage.out
fi

# Clean up
$COMPOSE_CMD -f docker-compose.test.yml down

exit $test_exit
