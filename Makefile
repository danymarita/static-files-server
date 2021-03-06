.PHONY: default help build package tag push helm helm-migration run test clean

SHELL         = /bin/bash
APP_NAME      = static-files-server
VERSION      := $(shell git describe --always --tags)
GIT_COMMIT    = $(shell git rev-parse HEAD)
GIT_DIRTY     = $(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE    = $(shell date '+%Y-%m-%d-%H:%M:%S')
REGISTRY_URL  = "danympradana"

default: help

help:
	@echo 'Management commands for ${APP_NAME}:'
	@echo
	@echo 'Usage:'
	@echo '    make build                 Compile the project.'
	@echo '    make package               Build final Docker image with just the Go binary inside.'
	@echo '    make tag                   Tag image created by package with latest, git commit and version.'
	@echo '    make push                  Push tagged images to registry.'
	@echo '    make run ARGS=             Run with supplied arguments.'
	@echo '    make test                  Run tests on a compiled project.'
	@echo '    make test-cover            Run tests with goveralls.'
	@echo '    make clean                 Clean the directory tree.'

	@echo

build:
	@echo "Building ${APP_NAME} ${VERSION}"
	env GOOS=linux GOARCH=amd64 go build -ldflags "-w -X github.com/danymarita/static-files-server/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/danymarita/static-files-server/version.Version=${VERSION} -X github.com/danymarita/static-files-server/version.Environment=${ENVIRONMENT} -X github.com/danymarita/static-files-server/version.BuildDate=${BUILD_DATE}" -o bin/${APP_NAME}

package:
	@echo "Building image ${APP_NAME} ${VERSION} ${GIT_COMMIT}"
	docker build --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=${GIT_COMMIT}${GIT_DIRTY} -t ${REGISTRY_URL}/${APP_NAME}:local .

tag: package
	@echo "Tagging: latest ${VERSION} ${GIT_COMMIT}"
	docker tag ${REGISTRY_URL}/${APP_NAME}:local ${REGISTRY_URL}/${APP_NAME}:${GIT_COMMIT}
	docker tag ${REGISTRY_URL}/${APP_NAME}:local ${REGISTRY_URL}/${APP_NAME}:${VERSION}
	docker tag ${REGISTRY_URL}/${APP_NAME}:local ${REGISTRY_URL}/${APP_NAME}:latest

push: tag
	@echo "Pushing Docker image to registry: latest ${VERSION} ${GIT_COMMIT}"
	docker push ${REGISTRY_URL}/${APP_NAME}:${GIT_COMMIT}
	docker push ${REGISTRY_URL}/${APP_NAME}:${VERSION}
	docker push ${REGISTRY_URL}/${APP_NAME}:latest

run: build
	@echo "Running ${APP_NAME} ${VERSION}"
	bin/${APP_NAME} ${ARGS}

test:
	@echo "Testing ${APP_NAME} ${VERSION}"
	go test ./...

test-cover:
	@echo "Testing ${APP_NAME} ${VERSION}"
	go get github.com/mattn/goveralls && go test -v -cover -race ./...

clean:
	@echo "Removing ${APP_NAME} ${VERSION}"
	@test ! -e bin/${APP_NAME} || rm bin/${APP_NAME}
