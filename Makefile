.DEFAULT_GOAL := upload-agent-binary

build:
		mkdir -p bin && \
    	cd cmd/agent && \
    	GOOS=linux GOARCH=amd64 go build -o ../../bin/mear-agent && \
    	cd ../cli && \
    	go build -o ../../bin/mear

upload-agent-binary: build
		AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin aws s3 cp --endpoint-url=http://minio:9000 ./bin/mear-agent s3://mear-bin