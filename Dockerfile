FROM golang:1.25.10-alpine AS builder

WORKDIR /messenger

RUN apk add --no-cache build-base libwebp-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o messenger .

FROM alpine:latest

WORKDIR /messenger

RUN apk add --no-cache libwebp

COPY --from=builder /messenger/messenger .

EXPOSE 8086

CMD ["./messenger"]