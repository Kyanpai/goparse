VERSION ?= $(shell git describe --tags --always)
GO = go
GOFMT = gofmt

VERBOSE			= 0
C				= $(if $(filter 1,$(VERBOSE)),,@) ## Conditional command display
M				= @echo "\033[0;35m▶\033[0m"

.PHONY: all
all: vendor

.PHONY: init_vendor
init_vendor:
	$(M) running mod init…
	$(GO) mod init

.PHONY: vendor
vendor:
	$(M) running mod vendor…
	$(GO) mod vendor

.PHONY: tidy
tidy:
	$(M) running mod tidy…
	$(GO) mod tidy

.PHONY: test
test:
	$(M) running go test…
	$(GO) test -cover -v ./...

.PHONY: fmt
fmt:
	$(M) running mod fmt…
	$(GO) fmt ./...
