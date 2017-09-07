
#!/bin/bash

# The only argument this script should ever be called with is '--verify-only'

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(dirname "${BASH_SOURCE}")/..
BINDIR=${REPO_ROOT}/bin

# Generate the internal clientset (pkg/client/clientset_generated/internalclientset)
${BINDIR}/client-gen "$@" \
	      --input-base "github.com/nikhita/kube-custom-controller/pkg/apis/" \
	      --input "githubnotif/" \
	      --clientset-path "github.com/nikhita/kube-custom-controller/pkg/client/" \
	      --clientset-name internalclientset \
	      --go-header-file "${GOPATH}/src/github.com/kubernetes/repo-infra/verify/boilerplate/boilerplate.go.txt"
# Generate the versioned clientset (pkg/client/clientset_generated/clientset)
${BINDIR}/client-gen "$@" \
		  --input-base "github.com/nikhita/kube-custom-controller/pkg/apis" \
		  --input "githubnotif/v1" \
	      --clientset-path "github.com/nikhita/kube-custom-controller/pkg/" \
	      --clientset-name "client" \
	      --go-header-file "${GOPATH}/src/github.com/kubernetes/repo-infra/verify/boilerplate/boilerplate.go.txt"
# generate lister
${BINDIR}/lister-gen "$@" \
		  --input-dirs="github.com/nikhita/kube-custom-controller/pkg/githubnotif" \
	      --input-dirs="github.com/nikhita/kube-custom-controller/pkg/apis/githubnotif/v1" \
	      --output-package "github.com/nikhita/kube-custom-controller/pkg/listers" \
	      --go-header-file "${GOPATH}/src/github.com/kubernetes/repo-infra/verify/boilerplate/boilerplate.go.txt"
# generate informer
${BINDIR}/informer-gen "$@" \
	      --go-header-file "${GOPATH}/src/github.com/kubernetes/repo-infra/verify/boilerplate/boilerplate.go.txt" \
	      --input-dirs "github.com/nikhita/kube-custom-controller/pkg/apis/githubnotif" \
	      --input-dirs "github.com/nikhita/kube-custom-controller/pkg/apis/githubnotif/v1" \
	      --internal-clientset-package "github.com/nikhita/kube-custom-controller/pkg/client/internalclientset" \
	      --versioned-clientset-package "github.com/nikhita/kube-custom-controller/pkg/client" \
	      --listers-package "github.com/nikhita/kube-custom-controller/pkg/listers" \
	      --output-package "github.com/nikhita/kube-custom-controller/pkg/informers"