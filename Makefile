CONTAINER_IMAGE?=line-notification
VERSION?=$(shell git tag --points-at HEAD)
PORT?=9090
APP?=goapp

run:
	go run cmd/main.go

clean:
	rm -f $(APP)
	rm -f .env

test: clean
	go test -v -cover ./...

env: test
	echo "CONTAINER_IMAGE=$(CONTAINER_IMAGE)" >> .env
	echo "VERSION=$(VERSION)" >> .env
	echo "PORT=$(PORT)" >> .env
	echo "APP=$(APP)" >> .env

build: env
	docker build . --no-cache -t $(CONTAINER_IMAGE):$(VERSION) -f build/Dockerfile

docker: build
	docker stop $(CONTAINER_IMAGE):$(VERSION) || true && docker rm $(CONTAINER_IMAGE):$(VERSION) || true \
	docker run --name $(CONTAINER_IMAGE) -p $(PORT):$(PORT) --rm \
		-e "PORT=$(PORT)" \
		$(CONTAINER_IMAGE):$(VERSION)