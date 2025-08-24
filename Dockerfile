FROM golang:alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /order-service

COPY app/go.mod /order-service/

RUN go mod download

COPY app/ /order-service/

RUN go build -o /order-service/build/main

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /order-service/build/main /app/
COPY /config.yaml /app/config

ENV CONFIG_PATH=/app/config

CMD [ "/app/main" ]