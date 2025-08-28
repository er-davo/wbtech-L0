FROM golang:alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /order-service

COPY app/go.mod /order-service/

RUN go mod download

COPY app/ /order-service/

RUN go build -o build/main cmd/main.go

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /order-service/build/main /app/
COPY /config.yaml /app/config
COPY /migrations /app/migrations
COPY app/public /app/public

ENV CONFIG_PATH=/app/config
ENV MIGRATION_DIR=/app/migrations

CMD [ "/app/main" ]