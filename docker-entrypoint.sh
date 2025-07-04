#!/bin/sh

# Exit immediately if a command exits with a non-zero status.
set -e

# Wait for the database to be ready
# The script will try to connect to the database every second for 30 seconds
# before giving up.
echo "Waiting for postgres..."
for i in $(seq 1 30); do
    if nc -z $DB_HOST $DB_PORT; then
        echo "Postgres is up - executing command"
        break
    fi
    sleep 1
done

# Construct the database URL from environment variables
DATABASE_URL="postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"

# Run the migrations
# The -path flag points to the directory with your .sql files.
# The -database flag provides the connection details.
# The up command applies all pending migrations.
echo "Running migrations..."
/app/migrate -path /app/migrations -database "$DATABASE_URL" up

# Execute the main application
# This will run the command specified as CMD in the Dockerfile
exec "$@"
