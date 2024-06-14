.PHONY: mod clean build install

mod:
	go mod download;

clean:
	rm -f json-log-parser

build: mod
	go build -o json-log-parser github.com/dbadham-fr/json-log-parser/cmd;

install: build
	go install cmd/json-log-parser.go;
