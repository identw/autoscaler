#!/usr/bin/env bash
set -xe

version=$1

gofmt cloudprovider/hetzner-identw/*.go 1>/dev/null
make docker-builder
make build-in-docker-arch-amd64
docker build -t docker.io/identw/cluster-autoscaler:${version} -f tmp/Dockerfile .
docker push docker.io/identw/cluster-autoscaler:${version}

