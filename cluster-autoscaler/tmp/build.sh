#!/usr/bin/env bash
set -xe

version=$1

gofmt cloudprovider/hetzner/*.go 1>/dev/null
make docker-builder
make build-in-docker
docker build -t docker.io/identw/cluster-autoscaler:${version} -f tmp/Dockerfile .
docker push docker.io/identw/cluster-autoscaler:${version}

