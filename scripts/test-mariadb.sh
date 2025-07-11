#!/bin/bash

set -e

CONTAINER_NAME="todo-mariadb-test"
IMAGE="mariadb:latest"
PORT=3306
ROOT_PASSWORD="password"
DATABASE="testdb"
TIMEOUT=60

# Clean up any existing container
if podman ps -a --format "{{.Names}}" | grep -q "^${CONTAINER_NAME}$"; then
    echo "Removing existing container..."
    podman stop $CONTAINER_NAME >/dev/null 2>&1 || true
    podman rm $CONTAINER_NAME >/dev/null 2>&1 || true
fi

# Start container
echo "Starting MariaDB container..."
podman run --rm -d \
  --name $CONTAINER_NAME \
  -e MARIADB_ROOT_PASSWORD=$ROOT_PASSWORD \
  -e MARIADB_DATABASE=$DATABASE \
  -p $PORT:3306 \
  $IMAGE

# Wait for DB to be ready
echo -n "Waiting for MariaDB to be ready (timeout: ${TIMEOUT}s)..."
start_time=$(date +%s)
while true; do
    # Try to connect and execute a simple query
    if podman exec $CONTAINER_NAME mariadb -h localhost -u root -p$ROOT_PASSWORD -e "SELECT 1" >/dev/null 2>&1; then
        break
    fi
    
    # Check timeout
    elapsed=$(( $(date +%s) - start_time ))
    if [ $elapsed -ge $TIMEOUT ]; then
        echo " timeout!"
        echo "Error: MariaDB did not become ready within ${TIMEOUT} seconds"
        echo "Container logs:"
        podman logs $CONTAINER_NAME
        echo "Trying to get error details..."
        podman exec $CONTAINER_NAME mariadb -u root -p$ROOT_PASSWORD -e "SHOW STATUS" || true
        exit 1
    fi
    sleep 1
    echo -n "."
done
echo " ready!"

# Verify database exists and create tables
echo -n "Setting up test database..."
podman exec -i $CONTAINER_NAME mariadb -u root -p$ROOT_PASSWORD < sql/setup_test_db.sql || {
    echo " failed!"
    echo "Error: Failed to setup test database"
    exit 1
}
echo " OK"

# Print connection info
echo ""
echo "MariaDB Test Container Info:"
echo "  Container Name: $CONTAINER_NAME"
echo "  Host: localhost"
echo "  Port: $PORT"
echo "  Database: $DATABASE"
echo "  Root Password: $ROOT_PASSWORD"
echo ""

# Run tests if requested
if [ "$1" = "--run-tests" ]; then
    echo "Running tests..."
    go test -v ./pkg/todo
    echo ""
fi

echo "To stop container:"
echo "  podman stop $CONTAINER_NAME"
