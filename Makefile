build:
	@go build -o ./bin/api 

run: build
	@./bin/api

test:
	@go test -v ./...

seed:
	go run scripts/seed.go

docker:
	@echo "building docker file"
	@docker build -t api .
	@echo "running api inside Docker container"
	@docker run --name aditya-reservation-api -p 3000:3000 -d go-api

