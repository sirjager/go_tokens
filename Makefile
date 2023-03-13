.PHONY: proto tidy test run

SERVICE_NAME=tokens
RPCS_DIR=../rpcs
PROTO_DIR=./pkg/proto

GO_RPC_DIR=$(RPCS_DIR)/$(SERVICE_NAME)/go
JS_RPC_DIR=$(RPCS_DIR)/$(SERVICE_NAME)/js
PY_RPC_DIR=$(RPCS_DIR)/$(SERVICE_NAME)/py

STATIK_OUT=./docs
SWAGGER_OUT=./docs/swagger

proto:
	- rm -rf $(GO_RPC_DIR) $(JS_RPC_DIR) $(PY_RPC_DIR)
	- mkdir -p $(GO_RPC_DIR) $(JS_RPC_DIR) $(PY_RPC_DIR)
	- rm -f $(SWAGGER_OUT)/*.swagger.json
	- rm -rf $(STATIK_OUT)/statik

	protoc \
	--proto_path=$(PROTO_DIR) --go_out=$(GO_RPC_DIR) --go_opt=paths=source_relative	\
	--go-grpc_out=$(GO_RPC_DIR) --go-grpc_opt=paths=source_relative	\
	--grpc-gateway_out=$(GO_RPC_DIR) --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=$(SWAGGER_OUT) --openapiv2_opt=allow_merge=true,merge_file_name=$(SERVICE_NAME) \
	--grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:$(JS_RPC_DIR) \
	$(PROTO_DIR)/*.proto
	statik -src=$(SWAGGER_OUT) -dest=$(STATIK_OUT)
	- python -m grpc_tools.protoc -I$(PROTO_DIR) \
	--python_out=$(PY_RPC_DIR) --grpc_python_out=$(PY_RPC_DIR) \
	$(PROTO_DIR)/*.proto


tidy:
	rm -f ./go.sum
	rm -rf ./vendor
	go get github.com/sirjager/rpcs@latest
	go mod tidy
	go mod vendor

test:
	go clean -testcache
	go test -v -cover ./... 

build:
	golint ./...
	go build -o ./dist/main ./cmd/main.go

run:
	go run ./cmd/main.go

setup:
	- go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	- go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	- go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	- go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	- go install github.com/goreleaser/goreleaser@latest
	- go install github.com/automation-co/husky@latest
	- go install golang.org/x/tools/cmd/cover@latest
	- go install github.com/ktr0731/evans@latest
	- go install github.com/rakyll/statik@latest
	- pip install -r ./requirements.txt
	- npm i -g protoc-gen-grpc-web @commitlint/cli @commitlint/config-conventional
	- husky init .


