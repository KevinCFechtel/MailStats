#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build  -o /go/bin/app -v ./cmd/imaparchive

#final stage
FROM alpine:latest
ENV CONFIG_FILE_PATH=/config/config.json
ENV RUN_AS_SERVICE=false
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
ENTRYPOINT ["/bin/sh", "-c", "/app -configFile=$CONFIG_FILE_PATH -runAsService=$RUN_AS_SERVICE"]
LABEL Name=imaparchive Version=1.0