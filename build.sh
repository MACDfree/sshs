#!/bin/bash -e

GOARCH=amd64 GOOS=linux go build -o bin/sshs-linux-amd64 ./