VERSION := $(shell git describe --tags --always || echo dev)

default: build

build: frontend
	go build \
		-ldflags "-s -w -X main.version=$(version)" \
		-mod=readonly \
		-trimpath

frontend: node_modules
	node ci/build.mjs

node_modules:
	npm ci

publish: frontend
	bash ./ci/build.sh
