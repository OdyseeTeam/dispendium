#!/bin/bash

 set -euo pipefail

 DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
 cd "$DIR"
 cd ".."
 DIR="$PWD"


echo "== Installing dependencies =="
GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports
GO111MODULE=off go get -u github.com/jteeuwen/go-bindata/...
go mod download


echo "== Checking dependencies =="
go mod verify
set -e


echo "== Compiling =="
export IMPORTPATH="github.com/lbryio/dispendium"
mkdir -p "$DIR/bin"
go generate -v
export VERSIONSHORT="${TRAVIS_COMMIT:-"$(git describe --tags --always --dirty)"}"
export VERSIONLONG="${TRAVIS_COMMIT:-"$(git describe --tags --always --dirty --long)"}"
export COMMITMSG="$(echo ${TRAVIS_COMMIT_MESSAGE:-"$(git show -s --format=%s)"} | tr -d '"' | head -n 1)"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o "./bin/dispendium" -asmflags -trimpath="$DIR" -ldflags "-X ${IMPORTPATH}/meta.version=${VERSIONSHORT} -X ${IMPORTPATH}/meta.versionLong=${VERSIONLONG} -X \"${IMPORTPATH}/meta.commitMsg=${COMMITMSG}\""
#echo "== Done building linux version $("$DIR/bin/dispendium" version) =="
#CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -v -o "./bin/dispendium-darwin" -asmflags -trimpath="$DIR" -ldflags "-X ${IMPORTPATH}/meta.version=${VERSIONSHORT} -X ${IMPORTPATH}/meta.versionLong=${VERSIONLONG} -X \"${IMPORTPATH}/meta.commitMsg=${COMMITMSG}\""
#echo "== Done building darwin version $("$DIR/bin/dispendium-darwin" version) =="
#CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -v -o "./bin/dispendium.exe" -asmflags -trimpath="$DIR" -ldflags "-X ${IMPORTPATH}/meta.version=${VERSIONSHORT} -X ${IMPORTPATH}/meta.versionLong=${VERSIONLONG} -X \"${IMPORTPATH}/meta.commitMsg=${COMMITMSG}\""
#echo "== Done building windows version $("$DIR/bin/dispendium.exe" version) =="

echo "$(git describe --tags --always --dirty)" > ./bin/dispendium.txt
chmod +x ./bin/dispendium
exit 0