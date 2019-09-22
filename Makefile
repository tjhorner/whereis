.PHONY: dist dist-win dist-macos dist-linux ensure-dist-dir build build-frontend install uninstall

GOBUILD=packr2 build -ldflags="-s -w"
INSTALLPATH=/usr/local/bin

ensure-dist-dir:
	@- mkdir -p dist

build-frontend:
	cd frontend && npm run build

dist-win: ensure-dist-dir
	# Build for Windows x64
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o dist/whereis-windows-amd64.exe *.go

dist-macos: ensure-dist-dir
	# Build for macOS x64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o dist/whereis-darwin-amd64 *.go

dist-linux: ensure-dist-dir
	# Build for Linux x64
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o dist/whereis-linux-amd64 *.go

dist: build-frontend dist-win dist-macos dist-linux clean

clean:
	packr2 clean

build: build-frontend
	@- mkdir -p bin
	$(GOBUILD) -o bin/whereis *.go
	make clean
	@- chmod +x bin/whereis

install: build
	mv bin/whereis $(INSTALLPATH)/whereis
	@- rm -rf bin
	@echo "whereis was installed to $(INSTALLPATH)/whereis. Run make uninstall to get rid of it, or just remove the binary yourself."

uninstall:
	rm $(INSTALLPATH)/whereis

run:
	@- go run *.go