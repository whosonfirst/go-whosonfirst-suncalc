CWD=$(shell pwd)
GOPATH := $(CWD)/vendor:$(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	rmdeps fmt bin

deps:   self
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-geojson"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-utils"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-httpony"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/suncalc-go"

vendor: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +

fmt:
	go fmt cmd/*.go
	go fmt *.go

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/wof-suncalc-server cmd/wof-suncalc-server.go
