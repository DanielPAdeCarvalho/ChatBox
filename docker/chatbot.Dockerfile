# Build stage
FROM golang:1.22.2-alpine as builder
WORKDIR /app

# Copy go.mod and go.sum for dependency management
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o main .

############## Final stage  ##############
FROM alpine:latest
WORKDIR /app

# Copy static, templates and binary from the builder stage
COPY --from=builder /app/main /app
COPY static /app/static
COPY templates /app/templates

# Create a non-root user and switch to it
RUN adduser -D myuser
USER myuser

CMD ["./main"]
