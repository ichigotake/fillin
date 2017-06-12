BIN = fillin

all: clean build

build: deps
	go build -o build/$(BIN) .

install: deps
	go install

cross: deps
	goxc -max-processors=8 -build-ldflags="" \
		-os="linux darwin freebsd netbsd windows" -arch="386 amd64 arm" -d . \
		-resources-include='README*' -n $(BIN)

deps:
	go get -d -v .

test: testdeps build
	rm -rf .test
	go test -v ./...

testdeps:
	go get -d -v -t .

lint: lintdeps build
	go vet
	golint -set_exit_status ./...

lintdeps:
	go get -d -v -t .
	go get -u github.com/golang/lint/golint

clean:
	rm -rf build .test snapshot debian
	go clean

.PHONY: build install cross deps test testdeps lint lintdeps clean
