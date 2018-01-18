GO_VERSION=1.8
GO_FILES=$(shell find . -type f -name "*.go")
BIN_DIR ?= bin
BRANCH := $(shell git branch | sed -n -e 's/^\* \(.*\)/\1/p' | sed -e 's/\//_/g')
TAG := ${BRANCH}-$(shell git rev-parse --short HEAD)
IMAGE_URL := gridx/backupd:${TAG}

DOCKER_RUN := docker run --rm -v "$$PWD:/go/src/github.com/grid-x/backupd" -w /go/src/github.com/grid-x/backupd golang:${GO_VERSION} bash -c


all: bin/server

test:
	go test -v $(shell glide nv)

lint:
	golint -set_exit_status $(shell glide nv)

bin/server: ${GO_FILES}
	go build -o ${BIN_DIR}/server github.com/grid-x/backupd/cmd/server

bin/server.linux: ${GO_FILES}
	GOOS=linux go build -o ${BIN_DIR}/server.linux github.com/grid-x/backupd/cmd/server

bin/backup.sh: ${GO_FILES}
	go build -o ${BIN_DIR}/backup.sh github.com/grid-x/backupd/cmd/backup.sh

docker: bin/server.linux
	docker build -t ${IMAGE_URL} -f Dockerfile .

push: docker
	docker push ${IMAGE_URL}

ci:
	docker run --rm -v "$$PWD:/go/src/github.com/grid-x/backupd" -w /go/src/github.com/grid-x/backupd golang:${GO_VERSION} bash -c 'curl https://glide.sh/get | sh && make bin/server.linux'

ci_test:
	docker run --rm -v "$$PWD:/go/src/github.com/grid-x/backupd" -w /go/src/github.com/grid-x/backupd golang:${GO_VERSION} bash -c 'curl https://glide.sh/get | sh && make test'

ci_lint:
	 ${DOCKER_RUN} 'curl https://glide.sh/get | sh && go get -u github.com/golang/lint/golint && make lint'
