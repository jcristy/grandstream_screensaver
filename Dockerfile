FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary
# CGO_ENABLED=0 is important for static linking, allowing use of scratch base image
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 go build -o grandstream_screensaver .

# Stage 2: Create the final lean image
FROM scratch

WORKDIR /app

# Expose the port your HTTP server will listen on
EXPOSE 8080

# Copy the built binary from the builder stage
COPY --from=builder /app/grandstream_screensaver ./grandstream_screensaver

# Run the application
ENTRYPOINT ["./grandstream_screensaver"]
