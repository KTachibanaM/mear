#!/usr/bin/env bash
set -e

AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3 rm --endpoint-url=http://minio:9000 s3://mear-destination/train.4k.avi

cat ./dev/host-partial-failure.json | jq -c | ./dist/mear-host_linux_amd64_v1/mear-host || true

train_4k_avi_actual_size=$(AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3api head-object --endpoint-url=http://minio:9000 --bucket mear-destination --key train.4k.avi --query ContentLength)
train_4k_avi_expected_size=40000000
if [ $train_4k_avi_actual_size -lt $train_4k_avi_expected_size ]; then
    echo train.4k.avi does not meet the expected size
    exit 1
else
    echo train.4k.avi meets the expected size
fi