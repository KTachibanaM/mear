#!/usr/bin/env bash
set -e

# Prepare mear-agent-testing container
docker rm -f mear-agent-testing
docker run -d --name=mear-agent-testing --platform=amd64 debian:bullseye sleep infinity
docker network connect mear-network mear-agent-testing

# Run those in the container (via cloud-init in a real cloud)
docker exec -it mear-agent-testing apt update
docker exec -it mear-agent-testing apt install -y curl
docker exec -it mear-agent-testing curl -sL http://minio-agent-binary:9000/bin/mear-agent -o /root/mear-agent
docker exec -it mear-agent-testing chmod +x /root/mear-agent
AGENT_ARGS_S3_TARGET_BASE64=$(base64 -w 0 ./scripts/demo-agent-args-s3-target.json)
docker exec -it mear-agent-testing /root/mear-agent $AGENT_ARGS_S3_TARGET_BASE64

# Tear down
docker rm -f mear-agent-testing
