#!/bin/bash
export GOOS=linux && go build -o hostgolang cmd/cmd.go
docker build -t "hostgo" .
docker tag "hostgo" registry.hostgolang.com/hostgo
docker push registry.hostgolang.com/hostgo