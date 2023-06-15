#!/usr/bin/env bash
set -e

AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3 cp --endpoint-url=http://minio-agent-binary:9000 ./bin/mear-agent s3://bin
