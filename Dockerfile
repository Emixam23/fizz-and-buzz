# Dockerfile References: https://docs.docker.com/engine/reference/builder/

############################################ MULTI STAGE BUILD PART 1 ##############################################

# Start from golang v1.17 base image
FROM golang:1.17 as builder

# Creating/Cd work directory
WORKDIR /app

# Copying sources
COPY . .

# Run go build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 /usr/local/go/bin/go build -ldflags="-s -w" -o /app/fizz-and-buzz-api .

############################################ MULTI STAGE BUILD PART 2 ##############################################

# Using alpine
FROM alpine

# Copying executable
COPY --from=builder /app/fizz-and-buzz-api .
COPY --from=builder /app/.env .
