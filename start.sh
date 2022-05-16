#!/bin/sh

set -e

echo "Run db migration..."
/app/migrate --version
cat /app/app.env
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "Start the app"
exec "$@"