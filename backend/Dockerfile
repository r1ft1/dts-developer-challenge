FROM golang:1.24-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY . . 

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o /main 

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["/main"]
