core-test:
	docker run --rm -ti --cap-add=SYS_PTRACE -e WORKDIR=$(shell pwd) -e GOPATH=$(GOPATH) -v $(shell pwd):$(shell pwd) \
	golang:1.11.8-stretch /bin/sh $(shell pwd)/script/core-test.sh