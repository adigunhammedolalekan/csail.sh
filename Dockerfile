FROM golang:alpine3.11
RUN apk add build-base
COPY . /app
RUN mkdir -p /var/kube
RUN touch /var/kube/config
COPY k8s-config /var/kube/config
WORKDIR /app
RUN go get ./...
ENV GOOS linux
RUN go build -o hostgolang cmd/cmd.go
ENTRYPOINT [ "./hostgolang" ]