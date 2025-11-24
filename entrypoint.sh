#!/bin/sh
set -e

echo "Running goose migrations..."

GOOSE_DSN="host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable"

goose -dir /app/migrations postgres "$GOOSE_DSN" up

echo "Starting app..."
/app/qa
