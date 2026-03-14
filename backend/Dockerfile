# Minimal Dockerfile (MVP)
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download || true
COPY . .
RUN go build -o /bin/api ./cmd/api

FROM alpine:3.19
WORKDIR /app
COPY --from=build /bin/api /usr/local/bin/api
ENV PORT=8080
EXPOSE 8080
CMD ["api"]
