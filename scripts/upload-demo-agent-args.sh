#!/usr/bin/env bash
set -e

AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3 cp --endpoint-url=http://minio-source:9000 ./scripts/demo-agent-args.json s3://src
