#!/bin/sh

set -e

echo "Run db migration..."
/app/migrate --version
echo "$DB_SOURCE"
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "Start the app"
exec "$@"