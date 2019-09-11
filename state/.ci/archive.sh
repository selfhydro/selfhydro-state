#! /bin/bash

set -e -x

apt-get update
apt-get install zip -y

VERSION=$(cat version/version)

cd ./selfhydro-state/state/

zip selfhydro-state-release.zip *.go go.*
mv selfhydro-state-release.zip ../../release
