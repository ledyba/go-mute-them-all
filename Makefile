.PHONY: run get clean

run: .bin/mute-them-all
	.bin/mute-them-all

get:
	go get -u "github.com/ChimeraCoder/anaconda"
	go get -u "github.com/Sirupsen/logrus"
	go get -u "github.com/fatih/color"

.bin/mute-them-all: $(shell find . -type f -name '*.go')
	gofmt -w .
	@mkdir -p .bin
	go generate github.com/ledyba/go-mute-them-all/...
	go build -o .bin/mute-them-all github.com/ledyba/go-mute-them-all/...


clean:
	go clean github.com/ledyba/go-mute-them-all/...
