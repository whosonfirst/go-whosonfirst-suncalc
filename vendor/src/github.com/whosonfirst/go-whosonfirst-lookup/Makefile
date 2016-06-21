CWD=$(shell pwd)
GOPATH := $(CWD)/vendor:$(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:	prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-lookup; then rm -rf src/github.com/whosonfirst/go-whosonfirst-lookup; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-lookup/providers
	cp lookup.go src/github.com/whosonfirst/go-whosonfirst-lookup/
	cp providers/*.go src/github.com/whosonfirst/go-whosonfirst-lookup/providers/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	rmdeps bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-geojson"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-utils"

bin:	self
	@GOPATH=$(GOPATH) go build -o bin/wof-lookup-fs cmd/wof-lookup-fs.go

fmt:
	go fmt providers/*.go
	go fmt *.go

