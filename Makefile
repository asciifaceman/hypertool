.PHONY: build-all build-linux build-windows build-mac

build-all: build-mac build-linux build-windows
	zip target/hypertool_mac
	zip target/hypertool_linux
	zip target/hypertool_win

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o target/hypertool_mac

build-linux:
	CC=gcc GOOS=linux GOARCH=amd64 go build -o target/hypertool_linux

build-windows:
	CGO_ENABLED=1 CC=x86_64-w64-mingw64-gcc GOOS=windows GOARCH=amd64 go build -o target/hypertool_win

