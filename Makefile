.DEFAULT_GOAL := upload-agent-binary

build:
		goreleaser release --clean --skip-validate --skip-publish

upload-agent-binary: build
		AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3 cp --endpoint-url=http://minio:9000 ./dist/mear-agent_linux_amd64_v1/mear-agent s3://mear-bin

test: upload-agent-binary
		./dev/download-demo-video.sh \
		&& ./dev/test-host-success.sh \
		&& ./dev/test-host-shared-destination-bucket.sh \
		&& ./dev/test-host-partial-failure.sh \
		&& ./dev/test-cli.sh
