#!/bin/bash

# Start the test database
docker-compose -f docker-compose.test.yml up -d mariadb-test

# Wait for database to be ready
echo "Waiting for database to be ready..."
while ! docker-compose -f docker-compose.test.yml exec mariadb-test mysqladmin ping -h localhost --silent; do
    sleep 1
done

# Run tests with the test database
export MARIADB_TEST_DSN="testuser:testpass@tcp(localhost:3306)/testdb"
go test ./pkg/todo/... -v

# Clean up
docker-compose -f docker-compose.test.yml down
