all:
	go build -o minecraftd
.PHONY: all

clean:
	go clean -i -modcache
.PHONY: clean
