#!/bin/bash
set -e
set -u
set -x
version=$(git describe --abbrev)

go build -o mango -ldflags "-X github.com/birchb1024/mango.Version=${version}" main.go
strip mango
GOOS=windows GOARCH=amd64 go build -o mango.exe -ldflags "-X github.com/birchb1024/mango.Version=${version}" main.go
GOOS=darwin GOARCH=amd64 go build -o mango_mac -ldflags "-X github.com/birchb1024/mango.Version=${version}" main.go

mkdir -p pkg
tar zcvf pkg/mango-"${version}".tgz ./mango ./mango.exe ./mango_mac ./README.md
