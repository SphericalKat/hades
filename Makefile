.phony: build
build:
	CGOENABLED=0 go build -o bin/main api/main.go

.phony: run
run:
	@~/.air -d -c .air.conf

.phony: docker_build
docker_build:
	docker image build -t gdgvit/hades-2.0 -f ./ops/images/go.Dockerfile .
