clean:
	rm -rf ./gen

setup:
	go mod tidy

swagger: clean
	mkdir ./gen
	swagger -q generate server --exclude-main -t ./gen -f swagger.yml -P zelic91/users/models.User

run:
	go run .