FROM golang:1.26.8-alpine3.24 AS build

WORKDIR /app

# Modules layer
COPY go.mod go.sum ./
#RUN go mod download

# Build layer
COPY ../../.. .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /myapp ./cmd/app

FROM alpine:3.24.1 AS run

COPY --from=build /myapp /myapp
COPY --from=build /app/.env /.env

EXPOSE 8080

CMD ["/myapp"]