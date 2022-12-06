build:
	go build -o bin/miya

run: build 
	bin/miya

test:
	go test -v ./...

testscover:
	go test -coverprofile tests_cover.out ./...
	go tool cover -html=tests_cover.out

totalcover:
	go test -coverprofile tests_cover.out ./...
	go tool cover -func tests_cover.out

clean:
	go clean
	rm bin/miya
