# Base stage for both dev and prod
FROM golang:1.25rc1-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .


# Development stage
FROM base AS dev
RUN go install github.com/air-verse/air@latest
EXPOSE 8800
CMD ["air", "-c", ".air.toml"]

# Production build stage
FROM base AS builder
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Production stage
FROM alpine:3.18 AS prod
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8800
CMD ["./main"]