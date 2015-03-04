#!/bin/bash -e

docker build --no-cache -t hjfreyer/base base
./rebuild.sh
