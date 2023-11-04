APP_SERVER_VERSION := 1
gobuild = go build -ldflags "-X main.buildVersion=$1 -X 'main.buildDate=$(shell date +'%Y/%m/%d %H:%M:%S')' -X main.buildCommit=$(shell git rev-parse HEAD)" -v -o $2 $3


.PHONY: build
build:
	$(call gobuild,${APP_SERVER_VERSION}, "cmd/cenarius/cenarius", "cmd/cenarius/main.go")

.PHONY: build_linux
build_linux:
	GOOS=linux GOARCH=amd64 $(call gobuild,${APP_SERVER_VERSION}, "cmd/cenarius/cenarius-linux", "cmd/cenarius/main.go")

.PHONY: build_windows
build_windows:
	GOOS=windows GOARCH=amd64 $(call gobuild,${APP_SERVER_VERSION}, "cmd/cenarius/cenarius.exe", "cmd/cenarius/main.go")

.PHONY: build_macos
build_macos:
	GOOS=darwin GOARCH=amd64 $(call gobuild,${APP_SERVER_VERSION}, "cmd/cenarius/cenarius.app", "cmd/cenarius/main.go")

.PHONY: test
test:
	go test -v -race -timeout 30s -v -covermode=atomic ./...

.DEFAULT_GOAL := build
