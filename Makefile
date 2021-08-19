binary = awsacc

.PHONY: release
release:
	GOOS=windows GOARCH=amd64 go build -o ./bin/$(binary)_windows_amd64
	GOOS=linux GOARCH=amd64 go build -o ./bin/$(binary)_linux_amd64
	GOOS=darwin GOARCH=amd64 go build -o ./bin/$(binary)_darwin_amd64

.PHONY: docker
docker:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/$(binary)_linux_docker_amd64 -v -ldflags '-w -extldflags '-static''
	docker build -t cbrgm/awsacc:latest .
