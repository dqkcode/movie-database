PROJECT_NAME=movies
BUILD_VERSION=$(shell cat VERSION)
DOCKER_IMAGE=$(PROJECT_NAME):$(BUILD_VERSION)
GO_BUILD_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on
GO_FILES=$(shell go list ./... | grep -v /vendor/)

.SILENT:


all: fmt vet test install build

build:
	$(GO_BUILD_ENV) go build -v -o $(PROJECT_NAME)-$(BUILD_VERSION).bin .

install:
	$(GO_BUILD_ENV) go install

vet:
	$(GO_BUILD_ENV) go vet $(GO_FILES)

fmt:
	$(GO_BUILD_ENV) go fmt $(GO_FILES)

test:
	$(GO_BUILD_ENV) go test $(GO_FILES) -cover -v

local: build
	./$(PROJECT_NAME)-$(BUILD_VERSION).bin
img:
	docker build -t movies -f deployment/docker/Dockerfile .
dcu: fmt vet install build img
	cd deployment/docker && docker-compose up
