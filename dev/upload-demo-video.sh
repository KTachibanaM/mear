#!/usr/bin/env bash
set -e

wget --no-clobber -P ./dev https://mear-dev.s3.us-west-2.amazonaws.com/train.4k.mp4
AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3 cp --endpoint-url=http://minio:9000 ./dev/train.4k.mp4 s3://mear-dev
