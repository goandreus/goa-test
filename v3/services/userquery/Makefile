#! /usr/bin/make
#
# Makefile for userquery micro-service
#

export GO111MODULE=on

SERVICE := userquery
VERSION := 1.0.0
PKG := "gitlab.com/wiserskills/v3/services/${SERVICE}"
PKG_LIST := $(shell go list ${PKG}/...  | grep -v ${PKG}/design/)

gen: ## Generates the Goa files
	mv _design design
	sed -i "" 's/\/\/goa.design\/plugins\/v3 v3.0.0/goa.design\/plugins\/v3 v3.0.0/g' go.mod
	sed -i "" 's/\/\/replace/replace/g' go.mod
	goa gen "${PKG}/design" && \
	goa example "${PKG}/design"
	mv design _design
	sed -i "" 's/goa.design\/plugins\/v3 v3.0.0/\/\/goa.design\/plugins\/v3 v3.0.0/g' go.mod
	sed -i "" 's/replace/\/\/replace/g' go.mod

test: ## Runs the unit tests
	@go test -short ${PKG_LIST}

build: ## Builds the micro service (server & client)
	@cd "$(GOPATH)/src/${PKG}" && \
	go build -a -o ./build/${SERVICE} ./cmd/${SERVICE} && \
	go build -a -o ./build/${SERVICE}-cli ./cmd/${SERVICE}-cli

run: ## Runs the built service locally
	@cd "$(GOPATH)/src/${PKG}/build" && \
	CLUSTER=fr0 EVENTBUS_CLUSTER=wiserskills EVENTBUS_URL=nats://localhost:4222 TOKEN_ACTIVE=true ./${SERVICE} 

docker: ## Builds the docker image locally (requires a Dockerfile)
	GOOS=linux GOARCH=amd64 go build -a -o ./build/${SERVICE} ./cmd/${SERVICE}
	@if [ -f Dockerfile ]; then\
		docker rm --force ${SERVICE};\
		docker rmi "wiserskills/${SERVICE}:${VERSION}";\
		docker build -t "wiserskills/${SERVICE}:${VERSION}" .;\
	else\
		echo "No Dockerfile";\
	fi

compose-up: ## Starts a local instance of the service with instances of the event bus and ArangoDB in docker
	@if [ -f compose.yml ]; then\
		docker-compose -p ws -f compose.yml up -d ;\
	else\
		echo "No compose file.";\
	fi

compose-down: ## Shutdowns the local containers
	@docker-compose -p ws -f compose.yml down

jsclient: ## Generates GRPC Web client
	@protoc --proto_path=gen/grpc/${SERVICE}/pb --js_out=import_style=commonjs,binary:js/ --grpc-web_out=import_style=commonjs,mode=grpcwebtext:js/ gen/grpc/${SERVICE}/pb/${SERVICE}.proto

help: ## Displays this help screen
	grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
