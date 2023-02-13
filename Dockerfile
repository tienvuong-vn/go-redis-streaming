FROM golang:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Specify the command to run when the container starts
CMD ["./main"]