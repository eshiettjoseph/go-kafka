FROM  golang:1.21-alpine AS build

RUN apk add --no-cache git

WORKDIR /tmp/producer

COPY src .

RUN go mod download

RUN go build -o ./out/producer .

FROM alpine:3.18

COPY --from=build /tmp/producer/out/producer /app/producer

EXPOSE 8080


CMD ["/app/producer"]