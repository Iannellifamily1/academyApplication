# Step 1: Use an official Golang runtime as the parent image
FROM golang:1.19.2-bullseye

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy all files from the current directory to the container
COPY . .

# Step 4: Download the Go module dependencies
RUN go mod download

# Step 5: Build the Go app
RUN go build -o /godocker

# Step 6: Expose the port that your Go app listens on (change this based on your application)
EXPOSE 3001

# Step 7: Command to run your application
CMD ["/godocker"]
