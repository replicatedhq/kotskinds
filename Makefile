.PHONY: all
all: schemas openapischema

.PHONY: generate
generate: controller-gen client-gen
	$(CONTROLLER_GEN) \
		object:headerFile=./hack/boilerplate.go.txt \
		paths=./apis/...
	$(CONTROLLER_GEN) \
		crd \
		+output:dir=./config/crds \
		paths=./apis/kots/v1beta1/... \
		paths=./apis/kots/v1beta2/...
	$(CLIENT_GEN) \
		--output-pkg=github.com/replicatedhq/kotskinds/client \
		--output-dir ./client \
		--clientset-name kotsclientset \
		--input-base github.com/replicatedhq/kotskinds/apis \
		--input kots/v1beta1,kots/v1beta2 \
		--go-header-file ./hack/boilerplate.go.txt


.PHONY: openapischema
openapischema: controller-gen
	$(CONTROLLER_GEN) \
		crd \
		+output:dir=./config/crds \
		paths=./apis/kots/v1beta1/... \
		paths=./apis/kots/v1beta2/...

.PHONY: schemas
schemas: fmt generate
	go build ${LDFLAGS} -o bin/schemagen ./schemagen
	./bin/schemagen --output-dir ./schemas

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: controller-gen
controller-gen:
ifeq (, $(shell which controller-gen))
	go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.19.0
CONTROLLER_GEN=$(shell go env GOPATH)/bin/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

.PHONY: client-gen
client-gen:
ifeq (, $(shell which client-gen))
	go install k8s.io/code-generator/cmd/client-gen@v0.34.1
CLIENT_GEN=$(shell go env GOPATH)/bin/client-gen
else
CLIENT_GEN=$(shell which client-gen)
endif

.PHONY: check-schemas
check-schemas: all
	@if [ -n "$$(git status --short)" ]; then \
		echo -e "\033[31mThe git repo is dirty :( Ensure all generated files using 'make all' are committed\033[0;m"; \
		git status --short; \
		exit 1; \
	fi
