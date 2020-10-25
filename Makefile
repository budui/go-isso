GOOS=linux
GOARCH=amd64

all: ilno

VERSION := 0.1.0
BUILD_DATE := `date +%FT%T%z`
LD_FLAGS := "-X 'wrong.wang/x/go-isso/version.Version=$(VERSION)' -X 'wrong.wang/x/go-isso/version.BuildTime=$(BUILD_DATE)'"


ilno: $(shell ag -l --go) bindata.go
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags $(LD_FLAGS) -o ilno

bindata.go: static/js/embed.min.js
	go-bindata -o bindata.go static/...

static/js/%.min.js: js/%.js $(shell find ./js) js/app/text/postbox.js js/app/text/comment_loader.js js/app/text/comment.js
	./node_modules/esbuild/bin/esbuild --bundle --sourcemap --outfile=$@ $<

js/app/text/%.js: js/app/text/%.jade
	./node_modules/.bin/pug --client --no-debug --name tt_$* --silent $<
	echo '' >> $@
	echo 'export default tt_$*' >> $@
