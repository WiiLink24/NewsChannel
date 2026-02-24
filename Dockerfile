FROM golang:alpine AS builder

# We assume only git is needed for all dependencies.
# openssl is already built-in.
RUN apk add -U --no-cache git

RUN adduser -D server
USER server
WORKDIR /home/server

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

# Get supercronic
RUN apk add -U --no-cache supercronic

RUN adduser -D server
USER server
WORKDIR /home/server

# Copy executable + country list
COPY --from=builder /home/server/app .
COPY countries.json .

COPY crontab .

CMD ["supercronic", "./crontab"]
