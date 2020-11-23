FROM golang:1.15-alpine as build

WORKDIR /app

COPY go.mod ./

RUN go mod download
RUN go mod verify

COPY . .

RUN go build -o main github.com/SERAGORN/siteparser/cmd

FROM alpine:latest

RUN apk update && apk add tzdata

COPY --from=build /app/main .

COPY config.yaml .

EXPOSE "3333"

ENTRYPOINT ["/main"]