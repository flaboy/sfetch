build:
	go build .

run:
	go run main.go

push:
	docker build . -t wanglei999/sfetch
	docker push wanglei999/sfetch