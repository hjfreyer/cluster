#!/bin/bash -x
REF=${1-master}
git subtree pull --prefix vagrant vagrant-origin $REF --squash
