PACKAGE_NAME := github.com/nikhita/kube-custom-controller

# A temporary directory to store generator executors in
BINDIR ?= bin
GOPATH ?= $HOME/gocode
HACK_DIR ?= hack

# A list of all types.go files in pkg/apis
TYPES_FILES = $(shell find pkg/apis -name types.go)

# This step pulls the Kubernetes repo so we can build generators in another
# target. Soon, github.com/kubernetes/kube-gen will be live meaning we don't
# need to pull the entirety of the k8s source code.
.get_deps:
	@echo "Grabbing dependencies..."
	@go get -d -u k8s.io/code-generator/cmd/... || true
	@go get -d github.com/kubernetes/repo-infra/... || true
	cd ${GOPATH}/src/k8s.io/code-generator;
	@touch $@

# Targets for building k8s code generators
#################################################
.generate_exes: .get_deps \
				$(BINDIR)/defaulter-gen \
                $(BINDIR)/deepcopy-gen \
                $(BINDIR)/conversion-gen \
                $(BINDIR)/client-gen \
                $(BINDIR)/lister-gen \
                $(BINDIR)/informer-gen
	touch $@

$(BINDIR)/defaulter-gen:
	go build -o $@ k8s.io/code-generator/cmd/defaulter-gen

$(BINDIR)/deepcopy-gen:
	go build -o $@ k8s.io/code-generator/cmd/deepcopy-gen

$(BINDIR)/conversion-gen:
	go build -o $@ k8s.io/code-generator/cmd/conversion-gen

$(BINDIR)/client-gen:
	go build -o $@ k8s.io/code-generator/cmd/client-gen

$(BINDIR)/lister-gen:
	go build -o $@ k8s.io/code-generator/cmd/lister-gen

$(BINDIR)/informer-gen:
	go build -o $@ k8s.io/code-generator/cmd/informer-gen
#################################################


# This target runs all required generators against our API types.
generate: .generate_exes $(TYPES_FILES)
	# Generate defaults
	$(BINDIR)/defaulter-gen \
		--v 1 --logtostderr \
		--go-header-file "$${GOPATH}/src/github.com/kubernetes/repo-infra/verify/boilerplate/boilerplate.go.txt" \
		--input-dirs "$(PACKAGE_NAME)/pkg/apis/github" \
		--input-dirs "$(PACKAGE_NAME)/pkg/apis/github/v1" \
		--extra-peer-dirs "$(PACKAGE_NAME)/pkg/apis/github" \
		--extra-peer-dirs "$(PACKAGE_NAME)/pkg/apis/github/v1" \
		--output-file-base "zz_generated.defaults"
	# Generate deep copies
	$(BINDIR)/deepcopy-gen \
		--v 1 --logtostderr \
		--go-header-file "$${GOPATH}/src/github.com/kubernetes/repo-infra/verify/boilerplate/boilerplate.go.txt" \
		--input-dirs "$(PACKAGE_NAME)/pkg/apis/github" \
		--input-dirs "$(PACKAGE_NAME)/pkg/apis/github/v1" \
		--output-file-base zz_generated.deepcopy
	# Generate conversions
	$(BINDIR)/conversion-gen \
		--v 1 --logtostderr \
		--go-header-file "$${GOPATH}/src/github.com/kubernetes/repo-infra/verify/boilerplate/boilerplate.go.txt" \
		--input-dirs "$(PACKAGE_NAME)/pkg/apis/github" \
		--input-dirs "$(PACKAGE_NAME)/pkg/apis/github/v1" \
		--output-file-base zz_generated.conversion
	# generate all pkg/client contents
	# Generate the internal clientset (pkg/client/clientset_generated/internalclientset)
	${BINDIR}/client-gen "$@" \
			--input-base "github.com/nikhita/kube-custom-controller/pkg/apis/" \
			--input "github/" \
			--clientset-path "github.com/nikhita/kube-custom-controller/pkg/client/" \
			--clientset-name internalclientset \
			--go-header-file "${GOPATH}/src/github.com/kubernetes/repo-infra/verify/boilerplate/boilerplate.go.txt"
	# Generate the versioned clientset (pkg/client/clientset_generated/clientset)
	${BINDIR}/client-gen "$@" \
			--input-base "github.com/nikhita/kube-custom-controller/pkg/apis" \
			--input "github/v1" \
			--clientset-path "github.com/nikhita/kube-custom-controller/pkg/" \
			--clientset-name "client" \
			--go-header-file "${GOPATH}/src/github.com/kubernetes/repo-infra/verify/boilerplate/boilerplate.go.txt"
	# generate lister
	${BINDIR}/lister-gen "$@" \
			--input-dirs="github.com/nikhita/kube-custom-controller/pkg/apis/github" \
			--input-dirs="github.com/nikhita/kube-custom-controller/pkg/apis/github/v1" \
			--output-package "github.com/nikhita/kube-custom-controller/pkg/listers" \
			--go-header-file "${GOPATH}/src/github.com/kubernetes/repo-infra/verify/boilerplate/boilerplate.go.txt"
	# generate informer
	${BINDIR}/informer-gen "$@" \
			--go-header-file "${GOPATH}/src/github.com/kubernetes/repo-infra/verify/boilerplate/boilerplate.go.txt" \
			--input-dirs "github.com/nikhita/kube-custom-controller/pkg/apis/github" \
			--input-dirs "github.com/nikhita/kube-custom-controller/pkg/apis/github/v1" \
			--internal-clientset-package "github.com/nikhita/kube-custom-controller/pkg/client/internalclientset" \
			--versioned-clientset-package "github.com/nikhita/kube-custom-controller/pkg/client" \
			--listers-package "github.com/nikhita/kube-custom-controller/pkg/listers" \
			--output-package "github.com/nikhita/kube-custom-controller/pkg/informers"