.PHONY: all
all: clean
	mkdir dist
	go build -o dist/tt ./cmd/tt

.PHONY: clean
clean:
	rm -rf dist