FROM golang:1.13-alpine as builder
RUN apk update && apk add git

WORKDIR /usr/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -o main .

FROM scratch

COPY --from=builder /usr/src/app .
CMD [ "/main" ]