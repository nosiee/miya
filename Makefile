build:
	go build -o bin/miya

run: build 
	bin/miya

test:
	go test -v ./...

clean:
	go clean
	rm bin/miya
