# Use an official Go runtime as a parent image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Download any dependencies
RUN go mod tidy

# Build the Go application
RUN go build -o main .

EXPOSE 8081 8082 8083

# Define the command to run your application using the provided ports
CMD ["./main", "8081", "8082", "8083"]