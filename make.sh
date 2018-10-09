#!/bin/bash

set -e

if [ -e agp.db ]; then
	echo "Exporting db schema..."
	sqlite3 agp.db .schema > agp_sample.db
fi

echo "Linting..."
gometalinter -e dashboard/assets.go

echo "Embedding assets..."
go generate dashboard/dashboard.go

rm awanagrandprix || true
go build .
