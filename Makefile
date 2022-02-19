buildRunDelete:
	go build .
	./melody-maker create
	rm ./melody-maker
