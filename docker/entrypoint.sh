#!/bin/sh
set -e

if [ -z "$POSTGRES_DSN" ]; then
  echo "POSTGRES_DSN is not set"
  exit 1
fi

attempt=0
max_attempts=10
until migrate -path /app/migrations -database "$POSTGRES_DSN" up; do
  attempt=$((attempt+1))
  if [ "$attempt" -ge "$max_attempts" ]; then
    echo "Migrations failed after $max_attempts attempts"
    exit 1
  fi
  echo "Migration failed, retrying in 3s..."
  sleep 3
done

echo "Starting API"
exec /app/health-map
