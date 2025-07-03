# Golang image as builder
FROM golang:1.24.4 AS builder

WORKDIR /app

ARG GITHUB_TOKEN
ENV GOPRIVATE=github.com/Platform48/*
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"


# Copy the mod and sum files
COPY metric_test/go.mod metric_test/go.sum ./
RUN go mod download

# Copy the source code
COPY metric_test/ ./

# Build 
RUN go build -o metricTest ./cmd/main.go

# Use minimal image for container
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/metricTest .

# Expose the port app hosted on
EXPOSE 8080

ENV FUNCTION_TARGET=MetricTest
# Command to run the binary
CMD ["./metricTest"]