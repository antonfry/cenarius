gobuild = go build -ldflags "-X main.buildVersion=$1 -X 'main.buildDate=$(shell date +'%Y/%m/%d %H:%M:%S')' -X main.buildCommit=$(shell git rev-parse HEAD)" -v -o $2 $3


.PHONY: build
build:
	$(call gobuild,${APP_SERVER_VERSION}, "cmd/cenarius/server", "cmd/cenarius/main.go")

.PHONY: test
test:
	go test -v -race -timeout 30s -v -covermode=atomic ./...

.DEFAULT_GOAL := build
