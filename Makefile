install:
	@echo 'Install the dependencies'
	go mod download

run:
	@echo 'Run [bin] ASOL CLI'
	cd bin && ./asol

build:
	@echo 'Build ASOL CLI'
	go build -o bin/asol

test:
	@echo 'Run tests'
	cd helper && go test -v helper_test.go helper.go