# Start from the latest Golang base image
FROM golang:latest

WORKDIR /app/common

# Copy the Common files
COPY ./common/go.mod ./common/go.sum ./
RUN go mod download

COPY ./common/application ./application
COPY ./common/domain ./domain
COPY ./common/infrastructure ./infrastructure
COPY ./common/web ./web

# Copy the feature-related files
WORKDIR /app/campaign

COPY ./campaign/go.mod ./campaign/go.sum ./
RUN go mod download

COPY ./campaign/application ./application
COPY ./campaign/domain ./domain
COPY ./campaign/infrastructure ./infrastructure
COPY ./campaign/web ./web

COPY ./campaign/docs ./docs
COPY ./campaign/main.go .

# Build the Go app
RUN go build -o main main.go

# Expose port 8080 to the outside world
EXPOSE 8081

# Command to run the executable
CMD ["./main"]
