FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build the app with static compilation
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/api

FROM gcr.io/distroless/base-debian11

WORKDIR /
COPY --from=builder /app/main .

# Expose the correct port for Cloud Run
EXPOSE 8080

CMD ["/main"]
