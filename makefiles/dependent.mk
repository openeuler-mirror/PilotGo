
GOLANGCILINT_VERSION ?= v1.59.1
GLOBAL_GOLANGCILINT := $(shell which golangci-lint)
GOBIN_GOLANGCILINT:= $(shell which ${GOBIN}/golangci-lint)


SWAG_VERSION ?= v1.16.3
GLOBAL_SWAGGERCI := $(shell which swag)
GOBIN_SWAGGERCI := $(shell which ${GOBIN}/swag)

.PHONY: golangci
golangci:
ifeq ($(shell $(GOBIN_GOLANGCILINT) version --format short), $(GOLANGCILINT_VERSION))
	@$(OK) golangci-lint is already installed
GOLANGCILINT=$(GOBIN_GOLANGCILINT)
else ifeq ($(shell $(GLOBAL_GOLANGCILINT) version --format short), $(GOLANGCILINT_VERSION))
	@$(OK) golangci-lint is already installed
GOLANGCILINT=$(GLOBAL_GOLANGCILINT)
else
	@{ \
	set -e ;\
	echo 'GOBIN path is: $(GOBIN)' ;\
	echo 'installing golangci-lint-$(GOLANGCILINT_VERSION)' ;\
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) $(GOLANGCILINT_VERSION) ;\
	echo 'Successfully installed' ;\
	}
GOLANGCILINT=$(GOBIN)/golangci-lint
endif

.PHONY: swagci
swagci:
ifeq ($(shell $(GOBIN_SWAGGERCI) --version | awk '{print $$3}'), $(SWAG_VERSION))
	@$(OK) swag is already installed ${GOBIN_SWAGGERCI}
SWAGCI=$(GOBIN_SWAGGERCI)
else ifeq ($(shell `$(GLOBAL_SWAGGERCI) --version` | awk '{print $$3}'), $(SWAG_VERSION))
	@$(OK) swag is already installed ${GLOBAL_SWAGGERCI}
SWAGCI=$(GLOBAL_SWAGGERCI)
else
	@{ \
	set -e ;\
	echo 'GOBIN path is: $(GOBIN)' ;\
	echo 'installing swag-$(SWAG_VERSION)' ;\
	go install -v github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION) ;\
	echo 'Successfully installed' ;\
	}
SWAGCI=$(GOBIN)/swag
endif