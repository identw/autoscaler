#!/usr/bin/env bash
set -xe

version=$1

gofmt cloudprovider/hetzner-identw/*.go 1>/dev/null
rm -fv ./cluster-autoscaler-arm64
rm -fv ./cluster-autoscaler-amd64
make docker-builder
make build-in-docker
docker build -t docker.io/identw/cluster-autoscaler:${version} -f tmp/Dockerfile .
docker push docker.io/identw/cluster-autoscaler:${version}

