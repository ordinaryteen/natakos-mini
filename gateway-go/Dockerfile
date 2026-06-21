# Use the official lightweight Go compilation workspace environment
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod ./
RUN go mod download

# Copy application source code files
COPY . .

# Compile the high-performance standalone binary executable target asset
RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway-go-binary main.go

# --- Final Runner Stage ---
# Use a razor-thin Alpine container image for our runtime execution layer
FROM alpine:3.19

WORKDIR /

# Pull the compiled binary file over from the heavy compilation stage
COPY --from=builder /gateway-go-binary /gateway-go-binary

EXPOSE 8085

CMD ["/gateway-go-binary"]