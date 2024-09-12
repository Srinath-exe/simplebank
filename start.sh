#!/bin/sh

set -e 

echo "Run DB Migrations"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "Start the server"
exec "$@"