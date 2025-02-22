# Use the latest Go 1.24 image
FROM golang:1.24

# Set working directory inside the container
WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all project files
COPY . .

# **Step 1: Build the training script**
RUN go build -o train_model train_model.go

# **Step 2: Run the training script inside the container to generate `model.json`**
RUN ./train_model

# **Step 3: Build the main API**
RUN go build -o energy-forecast server.go

# Expose API port
EXPOSE 5000

# Start the API
CMD ["./energy-forecast"]
