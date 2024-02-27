FROM  golang:1.19-alpine AS build

RUN apk add --no-cache git

WORKDIR /tmp/producer

COPY . .

RUN go mod download

RUN go build -o ./out/producer .

FROM alpine:3.18

COPY --from=build /tmp/producer/out/producer /app/producer

ENV KAFKA_SERVER_ADDRESS=""

EXPOSE 9092


CMD ["/app/producer"]