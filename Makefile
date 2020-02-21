SHELL := /bin/bash

PROTOC ?= $(shell which protoc)

PROTOS := pkg/api/schema/schema.proto pkg/api/outputs/outputs.proto pkg/api/version/version.proto pkg/api/inputs/inputs.proto
PROTO_URLS := https://raw.githubusercontent.com/falcosecurity/falco/feature/inputs-mvp/userspace/falco/schema.proto https://raw.githubusercontent.com/falcosecurity/falco/feature/inputs-mvp/userspace/falco/outputs.proto https://raw.githubusercontent.com/falcosecurity/falco/feature/inputs-mvp/userspace/falco/version.proto https://raw.githubusercontent.com/falcosecurity/falco/feature/inputs-mvp/userspace/falco/inputs.proto
PROTO_SHAS := a1f427c114b945d0880b55058862b74015d036aa722985ca6e5474ab4ed19f69 968957d993b97268223de9488f865ae062ccd8b60c0eb392237e46d4e278a912 0ac29f477e146a14504a34f91509668e44f4ee107d37b0be57e82847e8e89510 012f4d5aa5b935ef800591590674c499d7b6e492761bdeb16f105ef1cd99c440

PROTO_DIRS := $(dir ${PROTOS})
PROTO_DIRS_INCLUDES := $(patsubst %/, -I %, ${PROTO_DIRS})

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

.PHONY: clean
clean: ${PROTO_DIRS}
	@rm -rf $^
