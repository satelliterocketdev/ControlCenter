ifneq ($(wildcard /etc/ssl/certs/ca-certificates.crt),)
	export CURL_CA_BUNDLE=/etc/ssl/certs/ca-certificates.crt
endif

ifeq ($(shell which go),)
	export PATH := /usr/local/go/bin:$(PATH)
endif

ARCH=amd64

dependencies:
	go mod download
	go mod tidy

build-app:
	go build -o $(PWD)/bin/i69social ./src/api

build-static:
	cd ./src/api/static && npm i && npm run build

dev-static:
	cd ./src/api/static && npm run dev

run:
	go run $(PWD)/src/api/main.go
