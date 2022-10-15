build:
	go build -o server github.com/zer0tonin/crud_benchmark/server/src
	docker-compose build server

run: build
	docker-compose up server

down:
	docker-compose down --rmi=local --volumes

dev:
	docker-compose up dev

test:
	go test github.com/zer0tonin/crud_benchmark/server/src/...
