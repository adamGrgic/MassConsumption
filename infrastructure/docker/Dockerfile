# Build stage
FROM golang:1.23 AS build

# Install dependencies
RUN apt-get update && apt-get install -y curl && apt-get install unzip

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (to cache dependencies)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Install Bun (for TailwindCSS processing)
RUN curl -fsSL https://bun.sh/install | bash
ENV PATH="/root/.bun/bin:$PATH"

# Install TailwindCSS
RUN bun install && bun add tailwindcss

# Compile TailwindCSS
RUN bunx tailwindcss -i ./client/src/main.css -o ./static/css/styles.css

# Build the Go application
RUN go build -ldflags="-s -w" -o goth-todo main.go

# Runtime stage (smaller final image)
FROM ubuntu:22.04
WORKDIR /root/

# Install runtime dependencies (if needed)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the compiled Go binary
COPY --from=build /app/goth-todo .

# Copy static assets (Tailwind CSS, JS, etc.)
COPY --from=build /app/static ./static

# Expose application port
EXPOSE 8080

# Start the server
CMD ["./goth-todo"]
