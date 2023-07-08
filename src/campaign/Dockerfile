# Start from the latest Golang base image
FROM golang:latest

# Copy the Common files
WORKDIR /common

COPY go.mod go.sum ./
RUN go mod download

COPY ../common/application ./application
COPY ../common/domain ./domain
COPY ../common/infrastructure ./infrastructure
COPY ../common/web ./web

# Copy the feature-related files
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY application ./application
COPY domain ./domain
COPY infrastructure ./infrastructure
COPY web ./web

COPY docs ./docs
COPY startup.go .

# Build the Go app
RUN go build -o main main.go

# Expose port 8080 to the outside world
EXPOSE 8081

# Command to run the executable
CMD ["./main"]
