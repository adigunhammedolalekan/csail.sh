FROM golang:alpine3.11
RUN apk add build-base
COPY . /app
WORKDIR /app

RUN go get ./...
ENV GOOS linux
RUN go build -o hostgolang cmd/cmd.go
ENTRYPOINT [ "./hostgolang" ]