FROM golang:1.25.1-alpine AS gocart_dev
RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
WORKDIR /app
COPY . .
RUN go mod download
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]


FROM gocart_dev AS build_stage
RUN swag init -g cmd/api/main.go -o docs
RUN CGO_ENABLED=0 GOOS=linux go build \
  -tags stage \
  -ldflags="-s -w" \
  -o /app/api \
  ./cmd/api

FROM gocart_dev AS build_prod
RUN CGO_ENABLED=0 GOOS=linux go build \
  -tags prod \
  -ldflags="-s -w" \
  -o /app/api \
  ./cmd/api


FROM alpine:latest AS gocart_stage
WORKDIR /app
COPY --from=build_stage /app/api .
COPY --from=build_stage /app/docs ./docs
RUN addgroup -g 1000 appuser && \
  adduser -D -u 1000 -G appuser appuser
RUN chown -R appuser:appuser /app
USER appuser
EXPOSE 8080
CMD ["./api"]

FROM alpine:latest AS gocart_prod
WORKDIR /app
COPY --from=build_prod /app/api .
RUN addgroup -g 1000 appuser && \
  adduser -D -u 1000 -G appuser appuser
RUN chown -R appuser:appuser /app
USER appuser
EXPOSE 8080
CMD ["./api"]