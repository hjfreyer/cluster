#!/bin/bash -e

function build {
    docker build -t "hjfreyer/$1" $1
}

build transmission
build nginx
