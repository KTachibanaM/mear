#!/usr/bin/env bash
set -e

docker rm -f mear-agent-testing
docker run -d --name=mear-agent-testing --platform=amd64 debian:bullseye sleep infinity
docker network connect mear-network mear-agent-testing
docker cp ./bin/mear-agent mear-agent-testing:/root
docker cp ./scripts/demo.json mear-agent-testing:/root
docker exec -it mear-agent-testing /root/mear-agent /root/demo.json
docker rm -f mear-agent-testing
