GOOS=linux
GOARCH=amd64

.PHONY: build

VERSION := 0.1.0
BUILD_DATE := `date +%FT%T%z`
LD_FLAGS := "-X 'github.com/budui/go-isso/isso.Version=$(VERSION)' -X 'github.com/budui/go-isso/isso.BuildTime=$(BUILD_DATE)'"


build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags $(LD_FLAGS) .

init:
	git config core.hooksPath .githooks
