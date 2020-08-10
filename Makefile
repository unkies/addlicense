.PHONY: all
all: bin/addlicense

bin/addlicense:
	go build -o bin/addlicense ./cmd/addlicense/... 

.PHONY: clean
clean:
	rm -rf bin

.PHONY: dogfood
dogfood: bin/addlicense
	./bin/addlicense --license LICENSE --ignore testdata .

.PHONY: test
test:
	go test -v ./...