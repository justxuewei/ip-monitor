GO ?= go
BINARY ?= ipmonitor
PKG ?= ./cmd
VERSION_FILE ?= VERSION
VERSION ?= $(strip $(shell cat $(VERSION_FILE) 2>/dev/null || echo dev))
LDFLAGS ?= -X main.version=$(VERSION)
RELEASE_VERSION := $(word 2,$(MAKECMDGOALS))
RELEASE_EXTRA_ARGS := $(wordlist 3,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

ifneq ($(filter release,$(MAKECMDGOALS)),)
ifneq ($(RELEASE_VERSION),)
.PHONY: $(RELEASE_VERSION)
$(RELEASE_VERSION):
	@:
endif
endif

.PHONY: build release
build:
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY) $(PKG)

release:
	@if [ -z "$(RELEASE_VERSION)" ] || [ -n "$(RELEASE_EXTRA_ARGS)" ]; then \
		echo "Usage: make release <version>"; \
		exit 1; \
	fi
	@if [ "$$(git branch --show-current)" != "main" ]; then \
		echo "Cannot release: current branch is not main."; \
		exit 1; \
	fi
	@if ! git diff --quiet || ! git diff --cached --quiet || [ -n "$$(git ls-files --others --exclude-standard)" ]; then \
		echo "Cannot release: worktree has uncommitted changes."; \
		git status --short; \
		exit 1; \
	fi
	@if [ "$$(cat $(VERSION_FILE))" = "$(RELEASE_VERSION)" ]; then \
		echo "$(VERSION_FILE) is already $(RELEASE_VERSION)."; \
		exit 1; \
	fi
	@if git rev-parse -q --verify "refs/tags/$(RELEASE_VERSION)" >/dev/null; then \
		echo "Cannot release: tag $(RELEASE_VERSION) already exists."; \
		exit 1; \
	fi
	printf '%s\n' "$(RELEASE_VERSION)" > $(VERSION_FILE)
	git add $(VERSION_FILE)
	git commit -m "Release $(RELEASE_VERSION) version"
	git tag "$(RELEASE_VERSION)"
	git push origin main
	git push origin "$(RELEASE_VERSION)"
