#!/bin/sh
set -e
rm -rf manpages
mkdir manpages
go run ./cmd/photos2map/ man | gzip -c -9 >manpages/photos2map.1.gz
