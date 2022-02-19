buildTest:
	go build .
	./melody-maker create ./midi
	rm ./melody-maker
