#!/usr/bin/env bash
set -e

declare -A expected_outputs
expected_outputs["train.4k.avi"]=40000000
expected_outputs["train.720p.avi"]=5000000
expected_outputs["castle.720p.avi"]=1600000

for key in "${!expected_outputs[@]}"; do
    AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3 rm --endpoint-url=http://minio:9000 "s3://mear-destination/${key}"
done

cat ./dev/host-shared-destination-bucket.json | jq -c | ./dist/mear-host_linux_amd64_v1/mear-host

for key in "${!expected_outputs[@]}"; do
    actual_size=$(AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3api head-object --endpoint-url=http://minio:9000 --bucket mear-destination --key "${key}" --query ContentLength)
    expected_size=${expected_outputs[${key}]}
    if [ $actual_size -lt $expected_size ]; then
        echo ${key} does not meet the expected size
        exit 1
    else
        echo ${key} meets the expected size
    fi
done 
