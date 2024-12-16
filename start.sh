#!/bin/sh

# Source the .env file if it exists
if [ -f .env ]; then
  export $(cat .env | xargs)
fi

echo "Database Host: $DB_HOST"
echo "Database Port: $DB_PORT"
# Build the Go application
go build -o main .


# Start the application
./main