
.PHONY: all
all: build test images
	$(MAKE) reset-permission

.PHONY: build
build: go-build java-build
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
images: image-hello-world image-kafka-producer image-kafka-consumergroup

.PHONY: go-build
go-build: protc-go
	docker run -it --rm \
		--workdir /src \
		--volume $(PWD)/go-src:/src \
		--volume $(HOME)/go:/go \
		golang:1.12 \
		go build -v ./...

.PHONY: java-build
java-build: protc-java
	docker run -it --rm \
		--workdir /src \
		--volume $(PWD)/java-src:/src \
		--volume $(HOME)/.m2:/root/.m2 \
		maven \
		mvn compile

.PHONY: go-test
go-test: protc-go
	docker run -it --rm \
		--workdir /src \
		--volume $(PWD)/go-src:/src \
		--volume $(HOME)/go:/go \
		golang:1.12 \
		go test -v ./...

.PHONY: java-test
java-test: protc-java
	docker run -it --rm \
		--workdir /src \
		--volume $(PWD)/java-src:/src \
		--volume $(HOME)/.m2:/root/.m2 \
		maven \
		mvn test

.PHONY: image-hello-world
image-hello-world:
	docker build \
		--tag yarencheng/one-tree:hello-world-latest \
		--file docker/hello-world/Dockerfile \
		.

.PHONY: image-kafka-producer
image-kafka-producer:
	docker build \
		--tag yarencheng/one-tree:kafka-producer-latest \
		--file docker/kafka-producer/Dockerfile \
		.

.PHONY: image-kafka-consumergroup
image-kafka-consumergroup:
	docker build \
		--tag yarencheng/one-tree:kafka-consumergroup-latest \
		--file docker/kafka-consumergroup/Dockerfile \
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

.PHONY: protc-java
protc-java:
	mkdir -p go-src/protobuf
	docker run -it --rm \
		--volume $(PWD):/src \
		--workdir /src \
		yarencheng/protoc \
			--java_out=java-src/src/main/java \
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
	docker rmi yarencheng/one-tree:kafka-producer-latest || true
	docker rmi yarencheng/one-tree:kafka-consumergroup-latest || true

.PHONY: protoc-clean
protoc-clean: protoc-clean-go protoc-clean-java

.PHONY: protoc-clean-go
protoc-clean-go:
	rm -rvf go-src/protobuf

.PHONY: protoc-clean-java
protoc-clean-java:
	rm -rvf java-src/src/main/java/com/yarencheng/protobuf
