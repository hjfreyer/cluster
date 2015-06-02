#!/bin/bash -e

python2 mksecret.py sslcerts $HOME/keys/hjfreyer.com-cert/2015/assembled/ | kubectl create -f -

kubectl create -f k8s/services.yaml
kubectl create -f k8s/transmission.yaml
kubectl create -f k8s/nginx.yaml
