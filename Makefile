USER=stephenlb
NAME=showhook
REPO=$(USER)/$(NAME)

.PHONY: network start stop build sh

all: build

build:
	docker build . -t $(NAME) --build-arg COMMIT=$(shell git describe --always)

network:
	docker network create -d bridge --subnet 192.168.0.0/24 --gateway 192.168.0.1 $(NAME)

run: start
start:
	docker run -e PORT=$(PORT) -e NEIGHBORS=$(NEIGHBORS) -p $(PORT):$(PORT) --net=$(NAME) -d $(NAME)

stop:
	docker stop $(NAME)
	docker rm $(NAME)

sh:
	docker run -it $(NAME)
