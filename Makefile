IMAGE    ?= homeland/txwifi
NAME     ?= kinokochat
VERSION  ?= 1.0.4

all: build push

dev: dev_build dev_run

build-amd64:
	docker buildx build --load \
		--platform linux/amd64 \
		--tag homeland-social/txwifi:latest .

build-arm64:
	docker buildx build --load \
		--platform linux/arm64/v8 \
		--tag homeland-social/txwifi:latest .

build-arm32v7:
	docker buildx build --load \
		--platform linux/arm/v7 \
		--tag homeland-social/txwifi:latest .

build: build-amd64

clean:
	rm -rf txwifi.tar.gz

push:
	docker push $(IMAGE):arm32v7-$(VERSION)

dev_build:
	docker build -t $(IMAGE) ./dev/

dev_run:

txwifi.tar.gz:
	docker save homeland-social/txwifi:latest | gzip -9 > txwifi.tar.gz

gomod:
	GOPROXY="" go mod vendor && go mod tidy