FROM golang:1.23 AS builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build the app
COPY . .
RUN go build -o ./cmd/api .

# Final lightweight container
FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=builder /app/main .
CMD ["/main"]
