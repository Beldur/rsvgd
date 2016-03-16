# Information for Version string in Jenkins
ifndef WORKSPACE
	WORKSPACE := "."
else
	WORKSPACE := $(WORKSPACE)/checkout
endif

COMMIT_SHA := $(shell cd $(WORKSPACE) && git rev-parse --short HEAD)
COMMIT_NUMBER := $(shell cd $(WORKSPACE) && git rev-list --count HEAD)
CURRENT_DATETIME := $(shell date +"%Y%m%d.%H%M%S")

LDFLAGS="-s -X main.ServiceVersion=1.0.$(COMMIT_NUMBER)-$(COMMIT_SHA)-$(CURRENT_DATETIME)"

.PHONY: build test vet

all: deps test vet build

build:
	@mkdir -p _output
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o _output/rsvgd_linux_amd64 -ldflags $(LDFLAGS)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -o _output/rsvgd_darwin_amd64 -ldflags $(LDFLAGS)

deps:
	go get -u ./...

test:
	go test -cover -v ./...

vet:
	go vet ./...

build-container: build
	docker build -t rsvgd .
