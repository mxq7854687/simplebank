#!/bin/sh

set -e

echo "Run db migration..."
/app/migrate --version
echo "$DB_SOURCE"
nc -zv postgres 5432
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "Start the app"
exec "$@"