## simple makefile to log workflow
.PHONY: deploy

LOGLEVEL ?= 1
SWAGGER ?= 2
GOFLAGS ?= $(GOFLAGS:)
PWD = $(shell pwd)
export GOPATH =$(PWD)/_libs:$(PWD)
export GOBIN = $(PWD)/bin

pkg: deploy npm
	@tar zcvf deploy/`basename "$(CURDIR)"`_service.tar.gz bin/

deploy: pkgrm clearbin clean format build install conf swagger

quick: build conf
	@cd bin/; ./`basename "$(CURDIR)"` $(a)

run: deploy
	@cd bin/; ./`basename "$(CURDIR)"` $(a)

pkgrm:
	@rm -rf pkg/ > /dev/null
	@rm -rf _libs/pkg > /dev/null

clearbin:
	@rm -rf bin/* > /dev/null

build:
	@go build $(GOFLAGS) ./...

initdev:
	@echo "Creating DEV config copies from production versions"
	@echo "Please make sure to update the production config before deployment\n"
	@cp config/newApp/conf.json config/newApp/dev.json
	@cp config/newApp/api_conf.json config/newApp/api_dev.json
	@cp config/logger/logger.json config/logger/logger_dev.json
	@cp scripts/env.sh scripts/env_dev.sh
	@echo "***Configs creation successful.***"

conftest:
	@cp attributes_info.json /tmp/attributes_info.json

conf:
ifeq ($(FLORESTENV),DEV)
	@echo "***DEV Environment detected. Copying dev configs.***"
	@mkdir -p bin/conf
	@cp config/newApp/dev.json bin/conf/conf.json
	@cp config/newApp/api_dev.json bin/conf/api_conf.json
	@cp config/logger/logger_dev.json bin/conf/logger.json
	@cp attributes_info.json bin/attributes_info.json
	@cp src/common/migration/mysql/queries.sql bin/queries.sql
	@cp src/amenities/migrations/common/attribute/attribute_mappings.csv bin/attribute_mappings.csv
	@bash scripts/env_dev.sh
	@bash scripts/logger.sh $(LOGLEVEL) bin/conf/logger.json
	@mkdir -p bin/aggregator
	@cp aggregator/aggregator.js bin/aggregator/
	@cp aggregator/HttpService.js bin/aggregator/
	@cp aggregator/mapper.js bin/aggregator/
	@cp aggregator/package.json bin/aggregator/
	@cp aggregator/transform.js bin/aggregator/
	@cp aggregator/utils.js bin/aggregator/
else
	@mkdir -p bin/conf
	@cp config/newApp/conf.json bin/conf/conf.json
	@cp config/newApp/api_conf.json bin/conf/api_conf.json
	@cp config/logger/logger.json bin/conf/logger.json
	@cp attributes_info.json bin/attributes_info.json
	@cp src/common/migration/mysql/queries.sql bin/queries.sql
	@cp src/amenities/migrations/common/attribute/attribute_mappings.csv bin/attribute_mappings.csv
	@bash scripts/env.sh
	@bash scripts/logger.sh $(LOGLEVEL) bin/conf/logger.json
	@mkdir -p bin/aggregator
	@cp aggregator/aggregator.js bin/aggregator/
	@cp aggregator/HttpService.js bin/aggregator/
	@cp aggregator/mapper.js bin/aggregator/
	@cp aggregator/package.json bin/aggregator/
	@cp aggregator/transform.js bin/aggregator/
	@cp aggregator/utils.js bin/aggregator/
endif

npm:
	@cd bin/aggregator; npm install

run_back: deploy
	@killall `basename "$(CURDIR)"` >> /dev/null 2>&1; echo
	@cd bin/; ./`basename "$(CURDIR)"` &

## EXPERIMENTAL. Requires inotify-tools installed on system
## Please run sudo apt-get install inotify-tools before running make watch
## Only for development, do not run in production.
## 5 second delay will occur post save
watch: run_back
	@while true; do\
		inotifywait -e create,delete -r *.go . > /dev/null;\
		sleep 5;\
		echo "File changes detected. Rebuilding. Wait for a while.";\
		killall `basename "$(CURDIR)"` >> /dev/null 2>&1; echo;\
		make run_back;\
	done

install:
	@go get ./...

test: format clean install conftest
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
