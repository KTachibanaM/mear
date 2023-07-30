#!/usr/bin/env bash
set -e

rm -f ./dev/train.4k.cli.avi

./bin/mear -i ./dev/train.4k.mp4 --mear-stack dev --mear-agent-timeout 60 ./dev/train.4k.cli.avi