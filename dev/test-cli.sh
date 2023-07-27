#!/usr/bin/env bash
set -e

./bin/mear -i ./dev/train.4k.mp4 --mear-stack dev --mear-agent-timeout 60 output.avi