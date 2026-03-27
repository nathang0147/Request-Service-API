FROM golang:1.26.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV GOPROXY=off
ENV GOSUMDB=off
RUN mkdir -p .bin && go build -o .bin/request-service-api ./cmd/api

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/.bin/request-service-api /app/request-service-api

ENV PORT=8080

EXPOSE 8080

ENTRYPOINT ["/app/request-service-api"]
