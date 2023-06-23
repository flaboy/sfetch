build:
	go build .

push:
	docker build . -t wanglei999/sfetch
	docker push wanglei999/sfetch