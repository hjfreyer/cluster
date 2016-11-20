#!/bin/bash -ex

yell() { echo "$0: $*" >&2; }
die() { yell "$*"; exit 111; }
try() { "$@" || die "cannot $*"; }

BUILD_DIR=$1

[[ "$BUILD_DIR" == "" ]] && die "Usage: build.sh DOCKER_DIR"

POD=$(kubectl get pods --namespace kube-system -l k8s-app=kube-registry \
	      -o template --template '{{range .items}}{{.metadata.name}} {{.status.phase}}{{"\n"}}{{end}}' \
	     | grep Running | head -1 | cut -f1 -d' ')

trap 'kill $(jobs -pr)' SIGINT SIGTERM EXIT
kubectl port-forward --namespace kube-system $POD 31111:5000 &

sleep 1

REMOTE_NAME=localhost:31111/$(basename "$BUILD_DIR")
sudo docker build "$BUILD_DIR" -t "$REMOTE_NAME"
sudo docker push "$REMOTE_NAME"
