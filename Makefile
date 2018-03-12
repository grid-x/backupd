GO_FILES := $(shell find . -type f -name "*.go")
GO_BUILD := CGO_ENABLED=0 go build -ldflags="-w -s"
GO_TOOLS := gridx/golang-tools:master-839443d
GO_PROJECT := github.com/grid-x/backupd
DOCKER_RUN := docker run -it ${DOCKER_LINK} --rm -v $$PWD:/go/src/${GO_PROJECT} -w /go/src/${GO_PROJECT}
GO_RUN := ${DOCKER_RUN} ${GO_TOOLS} bash -c

BIN_DIR ?= bin

BRANCH := $(shell git branch | sed -n -e 's/^\* \(.*\)/\1/p' | sed -e 's/\//_/g')
TAG := ${BRANCH}-$(shell git rev-parse --short HEAD)
IMAGE_URL := gridx/backupd:${TAG}

all: bin/server

test:
	go test -v $(shell glide nv)

lint:
	golint -set_exit_status $(shell glide nv)

bin/server: ${GO_FILES}
	${GO_BUILD} -o ${BIN_DIR}/server github.com/grid-x/backupd/cmd/server

bin/server.linux: ${GO_FILES}
	GOOS=linux ${GO_BUILD} -o ${BIN_DIR}/server.linux github.com/grid-x/backupd/cmd/server

bin/backup.sh: ${GO_FILES}
	${GO_BUILD} -o ${BIN_DIR}/backup.sh github.com/grid-x/backupd/cmd/backup.sh

docker: bin/server.linux
	docker build -t ${IMAGE_URL} -f Dockerfile .

push: docker
	docker push ${IMAGE_URL}

ci_build:
	${GO_RUN} "make bin/server.linux"

ci_test:
	${GO_RUN} "make test"

ci_lint:
	${GO_RUN} "make lint"
