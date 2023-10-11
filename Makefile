gobuild = go build -ldflags "-X main.buildVersion=$1 -X 'main.buildDate=$(shell date +'%Y/%m/%d %H:%M:%S')' -X main.buildCommit=$(shell git rev-parse HEAD)" -v -o $2 $3


.PHONY: build
build:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/gstats/v1/gstats.proto
	$(call gobuild,${APP_AGENT_VERSION}, "cmd/agent/agent", "cmd/agent/main.go")
	$(call gobuild,${APP_SERVER_VERSION}, "cmd/server/server", "cmd/server/main.go")

.PHONY: test
test:
	go test -v -race -timeout 30s -v -covermode=atomic ./...

.DEFAULT_GOAL := build
