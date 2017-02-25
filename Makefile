# Copyright (c) 2017 Christian Saide <Supernomad>
# Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

PUSH_COVERAGE=""
BENCH_MAX_PROCS=1

setup_var_run:
	@echo "Creating /var/run/protond..."
	@mkdir /var/run/protond
	@chown $$SUDO_USER /var/run/protond

build_docker:
	@echo "Building test docker container..."
	@docker-compose build

compile:
	@echo "Compiling protond..."
	@go install github.com/Supernomad/protond

build_deps:
	@echo "Running go get to install build dependancies..."
	@go get -u golang.org/x/tools/cmd/cover
	@go get -u github.com/mattn/goveralls
	@go get -u github.com/golang/lint/golint
	@go get -u github.com/GeertJohan/fgt

deps:
	@echo "Running go get to install library dependancies..."
	@go get -t -v './...'

lint:
	@echo "Running fmt/vet/lint..."
	@fgt go fmt './...'
	@fgt go vet './...'
	@fgt golint './...'

race:
	@echo "Running unit tests with race checking enabled..."
	@go test -race './...'

bench:
	@echo "Running unit tests with benchmarking enabled..."
	@GOMAXPROCS=$(BENCH_MAX_PROCS) go test -bench . -benchmem './...'

unit:
	@echo "Running unit tests with benchmarking disabled..."
	@go test './...'

coverage:
	@echo "Running go cover..."
	@dist/coverage.sh $(PUSH_COVERAGE)

cleanup:
	@echo "Cleaning up..."
	@rm -f protond.pid

release:
	@echo "Generating release tar balls..."
	@go build github.com/Supernomad/protond
	@tar czf protond_$(VERSION)_linux_amd64.tar.gz protond LICENSE
	@rm -f protond

ci: build_deps deps lint compile unit coverage

dev: deps lint compile unit coverage cleanup

full: deps lint compile bench coverage cleanup
