
.PHONY: all
all: build test images
	$(MAKE) reset-permission

.PHONY: build
build: go-build
	$(MAKE) reset-permission

.PHONY: test
test: go-test
	$(MAKE) reset-permission

.PHONY: clean
clean: go-clean delete-images protoc-clean
	$(MAKE) reset-permission

.PHONY: reset-permission
reset-permission:
	sudo chown -R `id -u` .
	sudo chown -R `id -u` $(HOME)/go

.PHONY: images
images: hello-world-image

.PHONY: go-build
go-build: protc-go
	docker run -it --rm \
		--workdir /src \
		--volume $(PWD)/go-src:/src \
		--volume $(HOME)/go:/go \
		golang:1.12 \
		go build -v ./...

.PHONY: go-test
go-test:
	docker run -it --rm \
		--workdir /src \
		--volume $(PWD)/go-src:/src \
		--volume $(HOME)/go:/go \
		golang:1.12 \
		go test -v ./...

.PHONY: hello-world-image
hello-world-image:
	docker build \
		--tag yarencheng/one-tree:hello-world-latest \
		--file docker/hello-world/Dockerfile \
		.

.PHONY: protc-go
protc-go:
	mkdir -p go-src/protobuf
	docker run -it --rm \
		--volume $(PWD):/src \
		--workdir /src \
		yarencheng/protoc \
			--go_out=go-src \
			protobuf/*.proto

.PHONY: go-clean
go-clean:
	docker run -it --rm \
		--workdir /src \
		--volume $(PWD)/go-src:/src \
		--volume $(HOME)/go:/go \
		golang:1.12 \
		go clean -v ./...

.PHONY: delete-images
delete-images:
	docker rmi yarencheng/one-tree:hello-world-latest || true

.PHONY: protoc-clean
protoc-clean: protoc-clean-go

.PHONY: protoc-clean-go
protoc-clean-go:
	rm -rvf go-src/protobuf
