SHELL := /bin/bash

GO ?= $(shell which go)
PROTOC ?= $(shell which protoc)
MOCKGEN ?= $(shell which mockgen)

ifeq ($(MOCKGEN),)
# keep this in sync with the version of github.com/golang/mock in go.mod
$(shell $(GO) get github.com/golang/mock/mockgen@v1.4.3)
endif

TEST_FLAGS ?= -v -race

PROTOS := pkg/api/schema/schema.proto pkg/api/outputs/outputs.proto pkg/api/version/version.proto
PROTO_URLS := https://raw.githubusercontent.com/falcosecurity/falco/master/userspace/falco/schema.proto https://raw.githubusercontent.com/falcosecurity/falco/master/userspace/falco/outputs.proto https://raw.githubusercontent.com/falcosecurity/falco/master/userspace/falco/version.proto
PROTO_SHAS := 1adf7fbb2b92793a3cf490204314af7788ffd81655c4cedb40587a22db9c1915 5e3bdc564c4d38f7d70a8fe50e6022a733ed93197edff6b824a24c6a45fed6c3 fc470546c00273bafe20b53ab6b7e0784206b8f6f9a705df92994e89035a5dc4

PROTO_DIRS := $(dir ${PROTOS})
PROTO_DIRS_INCLUDES := $(patsubst %/, -I %, ${PROTO_DIRS})

all: protos mocks

.PHONY: protos
protos: ${PROTOS}

# $(1): the proto path
# $(2): the proto URL
# $(3): the proto SHA256
define download_rule
$(1):
	@rm -f $(1)
	@mkdir -p ${PROTO_DIRS}
	@curl --silent -Lo $(1) $(2)
	@echo $(3) $(1) | sha256sum -c
	@${PROTOC} ${PROTO_DIRS_INCLUDES} $(1) --go_out=plugins=grpc,paths=source_relative:$(dir $(1))
endef
$(foreach PROTO,$(PROTOS),\
	$(eval $(call download_rule,$(PROTO),$(firstword $(PROTO_URLS)),$(firstword $(PROTO_SHAS))))\
	$(eval PROTO_URLS := $(wordlist 2,$(words $(PROTO_URLS)),$(PROTO_URLS)))\
	$(eval PROTO_SHAS := $(wordlist 2,$(words $(PROTO_SHAS)),$(PROTO_SHAS)))\
)

MOCK_PROTOS := pkg/api/outputs/outputs.proto pkg/api/version/version.proto
MOCK_SYMBOLS := ServiceClient,Service_GetClient,Service_SubClient ServiceClient

# $(1): the proto path
# $(2): the mock directory
# $(3): the mock filename
# $(4): the symbols to mock
define generate_mock
$(2)/$(3): $(1) protos
	@mkdir -p $(2)
	$(MOCKGEN) $(shell cat $(1) | sed -n -e 's/^option go_package = "\(.*\)";/\1/p') $(4) > $(2)/$(3)
endef
$(foreach PROTO,$(MOCK_PROTOS),\
	$(eval $(call generate_mock,$(PROTO),$(dir $(PROTO))mocks,$(patsubst %.proto,%.go,$(notdir $(PROTO))),$(firstword $(MOCK_SYMBOLS))))\
	$(eval MOCK_SYMBOLS := $(wordlist 2,$(words $(MOCK_SYMBOLS)),$(MOCK_SYMBOLS)))\
)

MOCKS := $(join $(dir ${MOCK_PROTOS}),$(patsubst %.proto,mocks/%.go,$(notdir ${MOCK_PROTOS})))

.PHONY: mocks
mocks: protos ${MOCKS}

.PHONY: test
test: mocks
	@$(GO) vet ./...
	@$(GO) test ${TEST_FLAGS} ./...

.PHONY: clean
clean: ${PROTO_DIRS}
	@rm -rf $^
