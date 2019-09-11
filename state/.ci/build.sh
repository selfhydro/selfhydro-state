#! /bin/bash

set -e -x

apt-get update
apt-get install zip -y

VERSION=$(cat version/version)

cd ./selfhydro-state/state/
GOOS=linux go build -v -ldflags '-d -s -w' -o writeStateToDynamoDB writeStateToDynamoDB.go
chmod +x writeStateToDynamoDB

zip selfhydro-state-release.zip writeStateToDynamoDB
mv selfhydro-state-release.zip ../../release
