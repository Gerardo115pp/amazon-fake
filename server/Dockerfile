#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./...

#final stage
FROM alpine:3.13
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
WORKDIR /go/src/app
ENTRYPOINT /app
EXPOSE 5006
CMD [ "gserver" ]
