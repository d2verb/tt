.PHONY: all
all:
	go build ./cmd/tt

.PHONY: clean
clean:
	rm -f tt