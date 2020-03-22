#!/bin/bash
export GOOS=linux && go build -o hostgolang cmd/cmd.go
docker build -t "hostgo" .
docker tag "hostgo" registry.csail.app/hostgo
docker push registry.csail.app/hostgo