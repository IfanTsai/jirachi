.PHONY: clean

all: debug

debug:
	go build -o bin/jirachi main.go

release:
	go build -ldflags '-s -w -X github.com/IfanTsai/jirachi/repl.Release=true' -o bin/jirachi main.go

clean:
	rm -r bin/jirachi
