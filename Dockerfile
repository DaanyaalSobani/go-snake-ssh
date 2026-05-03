# syntax=docker/dockerfile:1

# ---- build stage ----
FROM golang:1.26-alpine AS build

WORKDIR /src

# Cache deps separately so code edits don't bust the dep layer
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/snake-ssh ./cmd/snake-ssh

# ---- runtime stage ----
FROM alpine:3.20

WORKDIR /app
COPY --from=build /out/snake-ssh /app/snake-ssh

# host_key is mounted in via docker-compose; the binary loads it from CWD
EXPOSE 2222

ENTRYPOINT ["/app/snake-ssh"]
