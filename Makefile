all: freeportd

freeportd: $(shell find . -name '*.go')
	cd cmd/freeportd && go build -o ../../freeportd

clean:
	rm -f freeportd

.PHONY: clean
