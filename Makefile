build:
	go build -o bin/miya

run: build 
	bin/miya --fname chip8-roms/Pong.ch8 --delay 1 --background-color 0x000000FF --pixel-color 0xFFFFFFFF

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
