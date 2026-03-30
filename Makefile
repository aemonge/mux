BINARY  := mux
PREFIX  ?= /usr/local

.PHONY: build test install clean

build:
	go build -o $(BINARY) .

test:
	go test ./...

install: build
	install -d $(PREFIX)/bin
	install -m 755 $(BINARY) $(PREFIX)/bin/$(BINARY)

clean:
	rm -f $(BINARY)
