#!/bin/bash -e

cd $GOPATH
echo $TRAVIS_TAG
#GOARCH=amd64 GOOS=linux go build -o bin/certstrap-${BUILD_TAG}-linux-amd64 ${REPO_PATH}