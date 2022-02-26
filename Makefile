swagger:
	rm -rf ./gen
	mkdir ./gen
	swagger generate server --exclude-main -t ./gen swagger.yaml