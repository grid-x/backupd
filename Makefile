GO_FILES=$(shell find . -type f -name "*.go")
BIN_DIR ?= bin

all: backup.sh

test:
	go test -v $(shell glide nv)

server: ${GO_FILES}
	go build -o ${BIN_DIR}/backup.sh github.com/jhedev/backupd/cmd/server

backup.sh: ${GO_FILES}
	go build -o ${BIN_DIR}/backup.sh github.com/jhedev/backupd/cmd/backup.sh
