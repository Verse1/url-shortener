FROM golang:alpine AS builder

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o index .

########

FROM alpine

RUN adduser -S -D -H -h /final appuser

USER appuser

COPY . /final

COPY --from=builder /app/index /final/

WORKDIR /final

EXPOSE 3000

CMD ["./index"]





