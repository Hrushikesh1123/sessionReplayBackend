# Use the official Go image as the base image
FROM golang:1.23.3-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source and the start script
COPY . .

# Ensure the start.sh script is executable
RUN chmod +x start.sh

# Expose port 3000 to the outside world
EXPOSE 8080

CMD ["./start.sh"]