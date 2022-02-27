buildRunDelete:
	go build main.go
	./melody-maker create
	rm ./melody-maker

build:
	go build main.go
	./melody-maker

test:
	go build main.go
	./test.sh
