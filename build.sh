#!/bin/bash -e

ORG_PATH="github.com/MACDfree"
REPO_PATH="${ORG_PATH}/sshs"
BUILD_TAG="${TRAVIS_TAG:-$(git describe --all | sed -e's/.*\///g')}"

export GOPATH=$HOME/gopath
export GOBIN=$GOPATH/bin

GOARCH=amd64 GOOS=linux go build -o bin/sshs-${BUILD_TAG}-linux-amd64 ${REPO_PATH}
GOARCH=amd64 GOOS=darwin go build -o bin/sshs-${BUILD_TAG}-darwin-amd64 ${REPO_PATH}
