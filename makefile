.PHONY: build clean deploy gomodgen run-local

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/users users/deliveries/lambda/main.go

clean:
	rm -rf ./bin ./vendor

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

run-local:
	PORT=8005 TABLE_NAME=example-users go run cmd/server/main.go
