test:
	go test ./...

test_watch:
	find . | entr bash -c "clear && make test"

install:
	go build
	mv disciplinarian /usr/local/bin/disciplinarian