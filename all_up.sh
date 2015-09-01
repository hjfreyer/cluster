#!/bin/bash -e

KEYS=/data/hjfreyer/keys/
python mksecret.py sslcerts $KEYS/hjfreyer.com-cert/2015/assembled/ | kubectl create -f -
python mksecret.py ddclient-creds $KEYS/google-domains/ | kubectl create -f -
python mksecret.py btcd-creds $KEYS/btcd/ |kubectl create -f -

kubectl create -f k8s/services.yaml
kubectl create -f k8s/transmission.yaml
kubectl create -f k8s/nginx.yaml
kubectl create -f k8s/ddclient.yaml
