#!/usr/bin/env bash
set -e

wget --no-clobber -P ./dev https://archive.org/download/MakeMine1948/MakeMine1948_256kb.rm
AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3 cp --endpoint-url=http://minio:9000 ./dev/MakeMine1948_256kb.rm s3://source
