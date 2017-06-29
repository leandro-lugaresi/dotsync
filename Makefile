IMPORT_PATH=github.com/leandro-lugaresi/dotsync
SOURCES=$(shell find . -name "*.go" | grep -v vendor/)
PACKAGES=$(shell go list ./... | grep -v vendor/)

default: test

init:
	go get -u github.com/kisielk/errcheck
	go get -u github.com/golang/dep/cmd/dep
	dep init
.PHONY: init

errcheck:
	errcheck -ignore Close $(PACKAGES)
.PHONY: errcheck

# testing
test:
	go test -race -v $(PACKAGES)
.PHONY: test

test-ci: linters-ci
	GORACE="halt_on_error=1" go test -race -v $(PACKAGES)
.PHONY: test-ci

test-color:
	GORACE="halt_on_error=1" go test -race -v $(PACKAGES) | \
		sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | \
		sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/'' | \
		sed ''/RUN/s//$$(printf "\033[34mRUN\033[0m")/''
.PHONY: test-color

# bound checking: http://klauspost-talks.appspot.com/2016/go17-compiler.slide#8
bound-check:
	go build -gcflags="-d=ssa/check_bce/debug=1" ./... 2>&1 | grep -v vendor

build:
	@for GOOS in linux darwin; do \
		echo "Building for OS: $$GOOS"; \
		CGO_ENABLED=0 GOARCH=amd64 GOOS=$$GOOS go build -a $(VERSION_FLAGS) --ldflags '-extldflags "-static"' -tags netgo -installsuffix netgo -o ./dist/exporter_$$GOOS $(IMPORT_PATH); \
	done