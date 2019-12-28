.phony: build
build:
	CGOENABLED=0 go build -o bin/main api/main.go

.phony: run
run:
	@~/.air -d -c .air.conf

.phony: docker_build
docker_build:
	docker image build -t docker.pkg.github.com/atechnohazard/hades/hades:latest -f ./ops/images/go.Dockerfile .

docker_push:
	docker push docker.pkg.github.com/atechnohazard/hades/hades:latest