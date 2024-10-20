FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ./out/urlshortener ./cmd/main

FROM scratch

COPY --from=builder /build/out/urlshortener /app/urlshortener

COPY --from=builder /build/app.env /app/app.env

EXPOSE 8080

CMD ["./app/urlshortener"]