#!/usr/bin/env bash

AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3 rb --endpoint-url=http://minio:9000 s3://mear-dev --force
docker rm -f mear-agent-testing