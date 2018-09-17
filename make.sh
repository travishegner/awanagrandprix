#!/bin/bash

set -e

gometalinter

rm awanagrandprix || true
go build .
