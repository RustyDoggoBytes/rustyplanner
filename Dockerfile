FROM golang:1.23-alpine AS build

ENV CGO_ENABLED=1

RUN apk add --no-cache \
    gcc \
    # Required for Alpine
    musl-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go install github.com/mattn/go-sqlite3

COPY . .

RUN go build  -ldflags='-s -w -extldflags "-static"' -o myapp .

FROM scratch

WORKDIR /app
ENV HOST="0.0.0.0"

COPY --from=build /app/myapp .

EXPOSE 8080

ENTRYPOINT ["/app/myapp"]