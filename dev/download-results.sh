#!/usr/bin/env bash
set -e

AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3 cp --endpoint-url=http://minio:9000 s3://mear-dev/output.mp4 ./dev
