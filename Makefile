#!make

LOCAL_BIN := $(CURDIR)/bin
FRONTEND_DIR := $(CURDIR)/frontend
PATH := $(PATH):$(LOCAL_BIN)

BUF := $(LOCAL_BIN)/buf
PROTOC_GEN_GO := $(LOCAL_BIN)/protoc-gen-go
PROTOC_GEN_CONNECT_GO := $(LOCAL_BIN)/protoc-gen-connect-go
PROTO_TS_INSTALL_STAMP := $(FRONTEND_DIR)/node_modules/.proto-install-stamp

BUF_VERSION := v1.59.0
PROTOC_GEN_GO_VERSION := v1.36.11
PROTOC_GEN_CONNECT_GO_VERSION := v1.19.1

$(BUF):
	mkdir -p $(LOCAL_BIN)
	cd $(LOCAL_BIN) && GOTOOLCHAIN=local GOBIN=$(LOCAL_BIN) go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)

$(PROTOC_GEN_GO):
	mkdir -p $(LOCAL_BIN)
	cd $(LOCAL_BIN) && GOTOOLCHAIN=local GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)

$(PROTOC_GEN_CONNECT_GO):
	mkdir -p $(LOCAL_BIN)
	cd $(LOCAL_BIN) && GOTOOLCHAIN=local GOBIN=$(LOCAL_BIN) go install connectrpc.com/connect/cmd/protoc-gen-connect-go@$(PROTOC_GEN_CONNECT_GO_VERSION)

$(PROTO_TS_INSTALL_STAMP): $(FRONTEND_DIR)/package.json $(FRONTEND_DIR)/pnpm-lock.yaml
	CI=true pnpm --dir $(FRONTEND_DIR) install --frozen-lockfile
	touch $(PROTO_TS_INSTALL_STAMP)

.PHONY: proto-generate
proto-generate: proto-generate-go proto-generate-ts

.PHONY: proto-generate-go
proto-generate-go: $(BUF) $(PROTOC_GEN_GO) $(PROTOC_GEN_CONNECT_GO)
	$(BUF) generate --config buf.yml --template buf.gen.go.yml

.PHONY: proto-generate-ts
proto-generate-ts: $(BUF) $(PROTO_TS_INSTALL_STAMP)
	$(BUF) generate --config buf.yml --template buf.gen.ts.yml

.PHONY: proto-lint
proto-lint: $(BUF)
	$(BUF) lint --config buf.yml proto
