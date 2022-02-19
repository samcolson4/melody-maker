buildTest:
	go build .
	./melody-maker create blah
	rm ./melody-maker
