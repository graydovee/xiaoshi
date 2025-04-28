REPO=graydovee/xiaoshi
TAG=v0.1.3
IMG=$(REPO):$(TAG)
MCP_DIR=mcp-server

.Phony: build
build: mcp-server
	go build -o bin/xiaoshi main.go

.Phony: docker-build
docker-build:
	docker build -t $(IMG) -t $(REPO):latest .

.Phony: docker-push
docker-push:
	docker push $(IMG)
	docker push $(REPO):latest

.Phony: mcp-server
mcp-server: mcpserver-terminal

.Phony: mcpserver-terminal
mcpserver-terminal:
	go build -o ${MCP_DIR}/terminal mcpserver/terminal/main.go