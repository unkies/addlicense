.PHONY: all
all: bin/addlicense

bin/addlicense:
	go build -o bin/addlicense ./cmd/addlicense/... 

.PHONY: clean
clean:
	rm -rf bin