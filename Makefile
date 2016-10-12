all:
	make clean
	mkdir bin
	go build -o bin/fetcher fetcher/fetcher.go
	go build -o bin/generator generator/generator.go

clean:
	rm -rf ./bin ./colors.txt
