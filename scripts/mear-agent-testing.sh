#!/usr/bin/env bash
set -e

docker rm -f mear-agent-testing
docker run -d --name=mear-agent-testing --platform=amd64 debian:bullseye sleep infinity
docker network connect mear-network mear-agent-testing
docker cp ./bin/mear-agent mear-agent-testing:/root
AGENT_ARGS_S3_TARGET_BASE64=$(base64 -w 0 ./scripts/demo-agent-args-s3-target.json)
docker exec -it mear-agent-testing /root/mear-agent ${AGENT_ARGS_S3_TARGET_BASE64}
docker rm -f mear-agent-testing
