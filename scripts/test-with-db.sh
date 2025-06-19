#!/bin/bash

# Clean up any existing containers
podman-compose -f compose.test.yml down

# Start the test database with replacement
podman-compose -f compose.test.yml up -d --force-recreate mariadb-test

# Wait for database to be ready (with timeout)
echo "Waiting for database to be ready..."
timeout=30
while ! podman-compose -f compose.test.yml exec mariadb-test mysqladmin ping -h localhost --silent; do
    sleep 1
    timeout=$((timeout-1))
    if [ $timeout -le 0 ]; then
        echo "Timeout waiting for database"
        podman-compose -f compose.test.yml logs mariadb-test
        exit 1
    fi
done

# Run tests with the test database
export MARIADB_TEST_DSN="testuser:testpass@tcp(localhost:3306)/testdb"
go test ./pkg/todo/... -v
test_exit=$?

# Clean up
podman-compose -f compose.test.yml down

exit $test_exit
