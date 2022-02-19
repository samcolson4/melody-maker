buildRunDelete:
	go build .
	./melody-maker create
	rm ./melody-maker

build:
	go build .
	./melody-maker
