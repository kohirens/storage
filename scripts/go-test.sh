#!/bin/sh

set -e

apk --no-progress --purge add git

# Set the region of the bucket.
export AWS_REGION="us-east-2"

go mod tidy
go test -v ./...
