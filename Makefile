GOOS=linux
GOARCH=amd64

.PHONY: build

VERSION := 0.1.0
BUILD_DATE := `date +%FT%T%z`
LD_FLAGS := "-X 'wrong.wang/x/go-isso/version.Version=$(VERSION)' -X 'wrong.wang/x/go-isso/version.BuildTime=$(BUILD_DATE)'"


build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags $(LD_FLAGS) .

