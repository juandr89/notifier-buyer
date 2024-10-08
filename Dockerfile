# Step 1: Base stage (for dependency installation)
FROM golang:1.22.5-alpine AS base

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Step 2: Test stage
FROM base AS test
COPY . .
RUN go test ./... -v -coverprofile cover.out

# Step 3: Build stage
FROM base AS build
COPY . .
RUN go build -o my-go-app main.go

# Step 4: Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/my-go-app .
COPY config.yaml /app/config.yaml
EXPOSE 8080
CMD ["./my-go-app"]

