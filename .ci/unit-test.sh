#! /bin/bash

set -ex

export GOPATH=$PWD
mkdir bin
mkdir src
cd ./selfhydro-state
go get
go test -cover | tee test_coverage.txt

mv test_coverage.txt ../coverage-results/.
