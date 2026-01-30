FROM golang:1.24.11-alpine3.23 AS builder

# We assume only git is needed for all dependencies.
# openssl is already built-in.
RUN apk add -U --no-cache git

WORKDIR /app

# Cache pulled dependencies if not updated.
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy necessary parts of the source into the builder
COPY *.go ./
COPY news news

# Build to name "app".
RUN go build -o app .

# Runner
FROM alpine:latest

WORKDIR /app

# Copy executable + country list
COPY --from=builder /app/app .
COPY countries.json .

# Setup cron
COPY crontab .
RUN chmod 0644 ./crontab
RUN crontab ./crontab

CMD ["crond", "-f"]
