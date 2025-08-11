BINARY := kvasx
ARCHES := mips mipsle arm64
DIST := dist

.PHONY: build test release clean

build:
	go build -o $(DIST)/$(BINARY) ./cmd/$(BINARY)

test:
	go test ./...

release: clean
	mkdir -p $(DIST)
	for arch in $(ARCHES); do \
		GOOS=linux GOARCH=$$arch go build -o $(DIST)/$(BINARY)-$$arch ./cmd/$(BINARY); \
	done

clean:
	rm -rf $(DIST)

