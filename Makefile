all: minecraftd
.PHONY: all

minecraftd:
	go build -o minecraftd

test:
	go test -v
.PHONY: test

clean:
	go clean -modcache
.PHONY: clean
