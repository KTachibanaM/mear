#!/usr/bin/env bash
set -e

rm -f ./dev/train.4k.cli.avi

./dist/mear-cli_linux_amd64_v1/mear -i ./dev/train.4k.mp4 --mear-stack dev --mear-agent-timeout 60 ./dev/train.4k.cli.avi