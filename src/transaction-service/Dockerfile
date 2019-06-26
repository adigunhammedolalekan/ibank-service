FROM alpine:3.2

RUN apk update && apk add --no-cache ca-certificates

ADD . /app
WORKDIR /app

ENTRYPOINT [ "/app/transaction-service" ]