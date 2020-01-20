#!/bin/bash
docker build -t "hostgo" .
docker tag "hostgo" registry.hostgolang.com/hostgo
docker push registry.hostgolang.com/hostgo