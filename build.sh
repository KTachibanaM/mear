#!/bin/bash
set -e

mkdir -p bin

cd cmd/agent
go build -o ../../bin/mear-agent
cd ../host
go build -o ../../bin/mear
