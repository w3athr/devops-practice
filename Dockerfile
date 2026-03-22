FROM golang:1.25-alpine@sha256:8e02eb337d9e0ea459e041f1ee5eece41cbb61f1d83e7d883a3e2fb4862063fa \ 
    AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/weather-app

FROM alpine:3.23.3@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659 \ 
    AS app
RUN apk add --no-cache zlib=1.3.2-r0
WORKDIR /app
RUN adduser -D appuser
USER appuser
COPY --from=builder /app/weather-app /app/weather-app
EXPOSE 8000
CMD ["/app/weather-app"]
