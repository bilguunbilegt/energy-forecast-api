# Use Go official image
FROM golang:1.21

# Set working directory
WORKDIR /app

# Copy files
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN go build -o energy-forecast server.go

# Expose port
EXPOSE 5000

# Run the API
CMD ["./energy-forecast"]
