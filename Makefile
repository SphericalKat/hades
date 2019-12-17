.phony: build
build:
	CGOENABLED=0 go build -o bin/main api/main.go

.phony: run
run:
	@~/.air -d -c .air.conf