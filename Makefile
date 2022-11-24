all:
	go build -o bin/miya

run:
	bin/miya

clean:
	go clean
	rm bin/miya
