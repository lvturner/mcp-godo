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

# Start streaming logs in background
podman logs -f $CONTAINER_NAME > >(while read line; do echo "[mariadb] $line"; done) 2>&1 &
LOG_PID=$!

# Wait for DB to be ready
echo -n "Waiting for MariaDB to be ready (timeout: ${TIMEOUT}s)..."
start_time=$(date +%s)
while true; do
    # Check if port is accessible from host (bypasses container networking)
    if nc -z localhost $PORT 2>/dev/null; then
        break
    fi
    
    # Check timeout
    elapsed=$(( $(date +%s) - start_time ))
    if [ $elapsed -ge $TIMEOUT ]; then
        echo " timeout!"
        kill $LOG_PID 2>/dev/null || true
        echo "Error: MariaDB did not become ready within ${TIMEOUT} seconds"
        echo "Checking host port ${PORT}..."
        netstat -tuln | grep $PORT || true
        exit 1
    fi
    sleep 1
    echo -n "."
done
echo " ready!"

# Stop log streaming
kill $LOG_PID 2>/dev/null || true

# Verify database exists by connecting from host
echo -n "Verifying test database..."
if ! mysql -h 127.0.0.1 -P $PORT -u root -p$ROOT_PASSWORD -e "USE $DATABASE" 2>/dev/null; then
    echo " failed!"
    echo "Error: Could not access database $DATABASE"
    exit 1
fi
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
