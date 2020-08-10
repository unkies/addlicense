.PHONY: all
all: bin/addlicense

bin/addlicense:
	go build -o bin/addlicense ./cmd/addlicense/... 

.PHONY: clean
clean:
	rm -rf bin

.PHONY: dogfood_add
dogfood_add: bin/addlicense
	./bin/addlicense add --license LICENSE --ignore testdata .

.PHONY: dogfood_remove
dogfood_remove: bin/addlicense
	./bin/addlicense remove --license LICENSE --ignore testdata .

.PHONY: test
test:
	go test -v ./...