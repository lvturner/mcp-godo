#!/bin/bash

set -e

CONTAINER_NAME="todo-mariadb-test"
IMAGE="mariadb:latest"
PORT=3306
ROOT_PASSWORD="password"
DATABASE="testdb"

# Start container
podman run --rm -d \
  --name $CONTAINER_NAME \
  -e MARIADB_ROOT_PASSWORD=$ROOT_PASSWORD \
  -e MARIADB_DATABASE=$DATABASE \
  -p $PORT:3306 \
  $IMAGE

# Wait for DB to be ready
echo "Waiting for MariaDB to be ready..."
podman exec $CONTAINER_NAME mysqladmin ping -h localhost -u root -p$ROOT_PASSWORD --wait=30

echo "MariaDB container is ready (name: $CONTAINER_NAME)"
echo "Run tests with:"
echo "  go test -v ./pkg/todo"
echo ""
echo "To stop container:"
echo "  podman stop $CONTAINER_NAME"
