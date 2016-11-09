## simple makefile to log workflow
.PHONY: deploy

LOGLEVEL ?= 1
SWAGGER ?= 2
GOFLAGS ?= $(GOFLAGS:)
PWD = $(shell pwd)
export GOPATH =$(PWD)/_libs:$(PWD)
export GOBIN = $(PWD)/bin
	
pkg: deploy
	@tar zcvf deploy/`basename "$(CURDIR)"`_service.tar.gz bin/

deploy: pkgrm clean format build install conf swagger

build:
	@echo $(GOPATH)
	@go build $(GOFLAGS) ./...

run: deploy
	@cd bin/; ./`basename "$(CURDIR)"` $(a)

pkgrm:
	@rm -rf pkg/ > /dev/null
	@rm -rf _libs/pkg > /dev/null

clearbin:
	@rm -rf bin/* > /dev/null

env:
ifeq ($(FLORESTENV),DEV)
	$(shell bash -c "source scripts/env_dev.sh")
else
	$(shell bash -c "source scripts/env.sh")
endif

initdev:
	@echo "Creating DEV config copies from production versions"
	@echo "Please make sure to update the production config before deployment\n"
	@cp scripts/env.sh scripts/env_dev.sh
	@echo "***Configs creation successful.***"

conf:
	@mkdir -p bin/conf
	@cp config/newApp/conf.json bin/conf/conf.json
	@cp config/newApp/api_conf.json bin/conf/api_conf.json
	@cp config/logger/logger.json bin/conf/logger.json
	@bash scripts/logger.sh $(LOGLEVEL) bin/conf/logger.json

install: 
	@go get ./...

test: format clean install env
	@ginkgo -r -v=true -cover=true ./src/test/

coverage: install
	@sh src/test/coverage.sh

bench: install
	@go test -run=NONE -bench=. $(GOFLAGS) ./...

clean:
	@go clean $(GOFLAGS) -i ./...
	
format:
	@go fmt $(GOFLAGS) ./...
	
codeanalysis:
	@go tool vet src/	
	
swagger:
ifneq ($(SWAGGER),2)
	@echo building swagger support
	@bash scripts/swagger.sh $(SWAGGER)
endif

## EOF
