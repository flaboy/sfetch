build:
	go build .

run:
	go run main.go

push:
	docker buildx build --platform linux/amd64 . -t wanglei999/sfetch
	docker push wanglei999/sfetch