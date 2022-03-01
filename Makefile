SRC_DIRS = cmd pkg
GOSEC_VERSION := v2.8.1

.PHONY: default
default: build

.PHONY: build
build:  bindata control-plane-operator roks-metrics
	go build -mod=vendor -o ./bin/ibm-roks github.com/openshift/ibm-roks-toolkit/cmd/ibm-roks

.PHONY: bindata
bindata:
	hack/update-generated-bindata.sh

.PHONY: verify-bindata
verify-bindata:
	hack/verify-generated-bindata.sh

.PHONY: verify-gofmt
verify-gofmt:
	@echo Verifying gofmt
	@gofmt -l -s $(SRC_DIRS)>.out 2>&1 || true
	@[ ! -s .out ] || \
	  (echo && echo "*** Please run 'make fmt' in order to fix the following:" && \
	  cat .out && echo && rm .out && false)
	@rm .out

.PHONY: verify-gosec
verify-gosec:
	@echo Verifying gosec
	@curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b /tmp ${GOSEC_VERSION}
	@/tmp/gosec -exclude G104,G401,G402,G501 ./...

.PHONY: verify
verify: verify-gofmt verify-gosec verify-bindata

# Build manager binary
.PHONY: control-plane-operator
control-plane-operator:
	go build -mod=vendor -o ./bin/control-plane-operator ./cmd/control-plane-operator/main.go

.PHONY: roks-metrics
roks-metrics:
	go build -mod=vendor -o ./bin/roks-metrics ./cmd/roks-metrics/main.go
	go build -mod=vendor -o ./bin/metrics-pusher ./cmd/metrics-pusher/main.go

.PHONY: test-render
test-render: build
	./bin/ibm-roks render --config cluster.yaml.example --output-dir manifests --pull-secret empty-pull-secret.txt

.PHONY: create-release
create-release:
	@hack/create-release.sh

.PHONY: test
test:
	./hack/test.sh
