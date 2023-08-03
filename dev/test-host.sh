#!/usr/bin/env bash
set -e

rm -f ./dev/train.4k.avi
rm -f ./dev/castle.720p.avi

AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3 rm --endpoint-url=http://minio:9000 s3://mear-destination/train.4k.avi
AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3 rm --endpoint-url=http://minio:9000 s3://mear-destination/castle.720p.avi

cat ./dev/host.json | jq -c | ./dist/mear-host_linux_amd64_v1/mear-host

AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3 cp --endpoint-url=http://minio:9000 s3://mear-destination/train.4k.avi ./dev
AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3 cp --endpoint-url=http://minio:9000 s3://mear-destination/castle.720p.avi ./dev
