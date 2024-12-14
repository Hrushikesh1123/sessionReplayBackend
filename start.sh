#!/bin/sh

# Source the .env file if it exists
if [ -f .env ]; then
  export $(cat .env | xargs)
fi

# Build the Go application
go build -o main .

# Start the application
./main