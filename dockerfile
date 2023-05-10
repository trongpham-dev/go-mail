FROM golang:1.16-alpine AS build

WORKDIR /app

COPY . .

RUN go build -o main .

FROM alpine:3.14

WORKDIR /app

COPY --from=build /app/main .

COPY .env /app/.env

CMD ["/app/main"]