REPO=graydovee/xiaoshi
TAG=v0.0.12
IMG=$(REPO):$(TAG)

.Phony: build-linux
build-linux:
	GOOS=linux GOARCH=arm64 go build -o bin/xiaoshi main.go

.Phony: docker-build
docker-build:
	docker build -t $(IMG) -t $(REPO):latest .

.Phony: docker-push
docker-push:
	docker push $(IMG)
	docker push $(REPO):latest