#! /bin/bash

set -ex

export GOPATH=$PWD
cd ./selfhydro-state
go test -cover | tee test_coverage.txt

mv test_coverage.txt ../coverage-results/.
