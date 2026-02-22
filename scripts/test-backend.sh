#!/bin/sh

set -e

apk --no-progress --purge add git

cd backend
go mod tidy
go test -v ./...
