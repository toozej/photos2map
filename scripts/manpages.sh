#!/bin/sh
set -e
rm -rf manpages
mkdir manpages
go run ./cmd/golang-starter/ man | gzip -c -9 >manpages/golang-starter.1.gz
