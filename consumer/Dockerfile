FROM  golang:1.21-alpine AS build

RUN apk add --no-cache git

WORKDIR /tmp/consumer

COPY src .

RUN go mod download

RUN go build -o ./out/consumer .

FROM alpine:3.18

COPY --from=build /tmp/consumer/out/consumer /app/consumer

EXPOSE 8080

CMD ["/app/consumer"]