#!/usr/bin/env bash
set -e

cat ./dev/train.4k.json | jq -c | ./bin/mear-host