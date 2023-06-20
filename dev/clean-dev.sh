#!/usr/bin/env bash
set -e

AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3 rm --endpoint-url=http://minio-source:9000 s3://src/MakeMine1948_256kb.rm
AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3 rm --endpoint-url=http://minio-destination:9000 s3://dst/output.mp4
AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3 rm --endpoint-url=http://minio-destination:9000 s3://dst/agent.log
