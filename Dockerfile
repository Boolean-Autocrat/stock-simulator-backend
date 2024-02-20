# Build stage
FROM golang:1.21-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go mod download

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o main main.go

RUN apk --no-cache add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.19
RUN apk update && apk add gcc g++ libc-dev librdkafka-dev pkgconf
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY .env .
COPY start.sh .
COPY wait-for-it.sh .
COPY templates /app/templates
COPY assets /app/assets
COPY courseCodes.json .
RUN chmod +x start.sh
RUN chmod +x wait-for-it.sh
COPY db/migrations ./migration

EXPOSE 3000
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]