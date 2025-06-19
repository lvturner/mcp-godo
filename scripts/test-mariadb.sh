#!/bin/bash

set -e

CONTAINER_NAME="todo-mariadb-test"
IMAGE="mariadb:latest"
PORT=3306
ROOT_PASSWORD="password"
DATABASE="testdb"
TIMEOUT=30

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
while ! podman exec $CONTAINER_NAME mysqladmin ping -h localhost -u root -p$ROOT_PASSWORD --silent >/dev/null 2>&1; do
    sleep 1
    elapsed=$(( $(date +%s) - start_time ))
    if [ $elapsed -ge $TIMEOUT ]; then
        echo " timeout!"
        echo "Error: MariaDB did not become ready within ${TIMEOUT} seconds"
        exit 1
    fi
    echo -n "."
done
echo " ready!"

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
