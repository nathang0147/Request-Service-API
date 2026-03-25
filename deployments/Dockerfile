FROM golang:1.26.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/.bin/request-service-api /app/request-service-api

ENV PORT=8080

EXPOSE 8080

ENTRYPOINT ["/app/request-service-api"]
