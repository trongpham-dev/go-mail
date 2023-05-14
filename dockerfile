FROM golang:1.19-alpine AS build

WORKDIR /app

COPY . .

RUN go build -o main .

FROM alpine:3.14

WORKDIR /app

COPY --from=build /app/main .

COPY .env /app/.env

RUN apk add --no-cache tzdata

CMD ["/app/main"]