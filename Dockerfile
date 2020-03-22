FROM alpine:3.2
RUN apk update && apk add --no-cache ca-certificates
RUN apk add build-base
RUN mkdir -p /var/kube
RUN touch /var/kube/config
# COPY k8s-config /var/kube/config
COPY local-kube-config /var/kube/config
ADD . /app
WORKDIR /app
RUN chmod +x /app/hostgolang
EXPOSE 4005 5005
ENTRYPOINT [ "/app/hostgolang" ]