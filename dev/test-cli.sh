#!/usr/bin/env bash
set -e

rm -f ./dev/train.4k.cli.avi

./dist/mear-cli_linux_amd64_v1/mear -i ./dev/train.4k.mp4 --mear-stack dev --mear-agent-timeout 60 ./dev/train.4k.cli.avi

train_4k_avi_actual_size=$(stat -c%s "./dev/train.4k.cli.avi")
train_4k_avi_expected_size=40000000
if [ $train_4k_avi_actual_size -lt $train_4k_avi_expected_size ]; then
    echo train.4k.cli.avi does not meet the expected size
    exit 1
else
    echo train.4k.cli.avi meets the expected size
fi
