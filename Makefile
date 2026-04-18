# ISC License
#
# Copyright (c) 2026 Shane
#
# Permission to use, copy, modify, and/or distribute this software for any
# purpose with or without fee is hereby granted, provided that the above
# copyright notice and this permission notice appear in all copies.
#
# THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
# WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
# MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
# ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
# WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
# ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
# OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
ARTIFACT_NAME := changelog-updater-action
WORKING_DIR := $(shell pwd)
GOLIC_VERSION  ?= v0.1.2

VERSION ?= v0.0.0
GITSHA ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GOOS     ?= $(shell go env GOOS)
GOARCH   ?= $(shell go env GOARCH)
LD_FLAGS := "-X 'main.Version=$(VERSION)' -X 'main.Gitsha=$(GITSHA)'"

ifndef NO_COLOR
YELLOW=\033[0;33m
CYAN=\033[1;36m
RED=\033[31m
# no color
NC=\033[0m
endif

.PHONY: clean
clean::
	rm -rf $(WORKING_DIR)/bin

.PHONY: build
build::
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/$(ARTIFACT_NAME)-$(GOOS)-$(GOARCH) \
	-ldflags $(LD_FLAGS) \
	./cmd/action
	chmod +x $(WORKING_DIR)/bin/$(ARTIFACT_NAME)-$(GOOS)-$(GOARCH)

.PHONY: test
test::
	go test -v -tags=all -parallel ${TESTPARALLELISM} -timeout 2h

.PHONY: lint-init
lint-init:
	brew install golangci-lint
	brew install gitleaks
	brew install yamllint

.PHONY: lint
lint: license
	goimports -w ./
	golangci-lint run
	yamllint .
	gitleaks detect . --no-git --verbose --config=.gitleaks.toml

.PHONY: license
license: build
	golic inject -c "2026 Shane" -t isc

.PHONY: license-dry
license-dry: build
	golic inject -c "2026 Shane" -t isc -d

define golic
	@go install github.com/Bugs5382/golic/cmd/golic@$(GOLIC_VERSION)
	$(GOBIN)/golic inject $1
endef