#!/bin/sh

set -e

echo "run db migration"
/app/migrate -path /app/migration -database "$MIGRATE_DB_SOURCE" -verbose up

exec /app/main

psCheck=$(ps -ef | grep -v grep | grep -c "main")
if [ $psCheck -eq 0 ]; then
    echo "app is not running"
else
    echo "app is already running"
fi