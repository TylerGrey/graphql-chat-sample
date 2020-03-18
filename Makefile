schema:
	go-bindata -ignore=\.go -pkg=schema -o=./schema/bindata.go ./schema/...

run:
	go run main.go

redis:
	docker run --name graphql-redis -d -p 6379:6379 redis

.PHONY: schema run redis