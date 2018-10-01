#!/bin/bash

set -e

echo "Linting..."
gometalinter -e dashboard/assets.go

echo "Embedding assets..."
go generate dashboard/dashboard.go

rm awanagrandprix || true
go build .
