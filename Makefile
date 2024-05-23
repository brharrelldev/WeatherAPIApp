
API_KEY="enter key here"
build:
	go build -o bin/weather-api

test:
	go test ./... -cover

run-debug:
	go build -gcflags "all=-N -l" -o ./bin/weather-api
	./bin/weather-api -apikey $(API_KEY) -config config.yaml


run:
	./bin/weather-api -apikey $(API_KEY) -config config.yaml

