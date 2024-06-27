# Created by OpenAi

# Use an official Golang runtime as a parent image

FROM golang:1.17-alpine

WORKDIR /go/src/app

# Copy the source code and environment file to the Working Directory inside the container
COPY gym-webapp/backend/ .
COPY gym-webapp/backend/.env .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
