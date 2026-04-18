build:
	go build -o bin/qed .

run: build
	./bin/qed

test:
	go test ./... -v