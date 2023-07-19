# Start from the latest Golang base image
FROM golang:latest

# Copy the Common files
WORKDIR /app/common

COPY ./common/go.mod ./common/go.sum ./
RUN go mod download

COPY ./common/application ./application
COPY ./common/domain ./domain
COPY ./common/infrastructure ./infrastructure
COPY ./common/web ./web

# Copy the feature-related files
WORKDIR /app/coin

COPY ./coin/go.mod ./coin/go.sum ./
RUN go mod download

COPY ./coin/application ./application
COPY ./coin/domain ./domain
COPY ./coin/infrastructure ./infrastructure
COPY ./coin/web ./web

COPY ./coin/docs ./docs
COPY ./coin/main.go .

# Build the Go app
RUN go build -o main main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
