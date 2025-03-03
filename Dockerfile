# Use the latest Go 1.24 image
FROM golang:1.24

# Set working directory inside the container
WORKDIR /app

# Install AWS CLI (optional, useful for debugging S3 access)
RUN apt-get update && apt-get install -y awscli && rm -rf /var/lib/apt/lists/*

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all project files
COPY . .

# **Step 1: Build the training script**
RUN go build -o train_model train_model.go

# **Step 2: Set AWS credentials (if using environment variables)**
# Ensure the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are set at runtime.

# **Step 3: Run the training script inside the container to generate `model.json`**
RUN AWS_REGION=us-east-1 ./train_model

# **Step 4: Build the main API**
RUN go build -o energy-forecast server.go

# Expose API port
EXPOSE 5000

# Start the API
CMD ["./energy-forecast"]
