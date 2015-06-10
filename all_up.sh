#!/bin/bash -e

python2 mksecret.py sslcerts /data/hjfreyer/keys/hjfreyer.com-cert/2015/assembled/ | kubectl create -f -
python2 mksecret.py ddclient-creds /data/hjfreyer/keys/google-domains/ | kubectl create -f -

kubectl create -f k8s/services.yaml
kubectl create -f k8s/transmission.yaml
kubectl create -f k8s/nginx.yaml