FROM  golang:1.19-alpine AS build

RUN apk add --no-cache git

WORKDIR /tmp/consumer

COPY . .

RUN go mod download

RUN go build -o ./out/consumer .

FROM alpine:3.18

COPY --from=build /tmp/consumer/out/consumer /app/consumer

ENV KAFKA_SERVER_ADDRESS=""
ENV DB_HOST=""
ENV DB_USER=""
ENV DB_PASSWORD=""
ENV DB_NAME=""
ENV DB_PORT=""
ENV DB_SSLMODE=""

EXPOSE 9092

CMD ["/app/consumer"]