IMAGE    ?= kinokochat/txwifi
NAME     ?= kinokochat
VERSION  ?= 1.0.4

all: build push

dev: dev_build dev_run

build:
	docker build -t $(IMAGE):latest -t $(IMAGE):arm32v7-$(VERSION) .

push:
	docker push $(IMAGE):arm32v7-$(VERSION)

dev_build:
	docker build -t $(IMAGE) ./dev/

dev_run:
	sudo docker run --rm -it --privileged --network=host \
                   -v $(CURDIR):/go/src/github.com/kinokochat/txwifi \
                   -w /go/src/github.com/kinokochat/txwifi \
                   --name=$(NAME) $(IMAGE):latest

gomod:
	GOPROXY="" go mod vendor && go mod tidy