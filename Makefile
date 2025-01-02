# check OS and set the correct path separator and commands
ifeq ($(OS),Windows_NT)
	FIXPATH = $(subst /,\,$1)
	RM = del /Q
	MKDIR = mkdir
else
	FIXPATH = $1
	RM = rm -Rf
	MKDIR = mkdir -p
endif

GOEXE:=$(shell go env GOEXE)
GOOS:=$(shell go env GOOS)
GOARCH:=$(shell go env GOARCH)

NAME:=ged-journal
CMDS:= ged-journal simulator

# BINNAME:=$(call FIXPATH,dist/$(CMDNAME)$(GOEXE))


DIRS = dist var

help:
	@echo "make tidy"
	@echo "       go fmt and go mod tidy"
	@echo "make run"
	@echo "       run the application"
	@echo "make test"
	@echo "       run the tests"
# @echo "make clean"
# @echo "       remove the dist directory"

$(DIRS):
	$(MKDIR) $@

.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v


.PHONY: audit
audit:
	@echo "running checks"
	go mod verify
	go vet ./...
	go list -m all
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...


.PHONY: test
test:
	go test ./...

.PHONY: run-collect
run-collect: $(DIRS)
	go run $(call FIXPATH,./cmd/ged-journal) collect -f var/collect-log.jsonl --nats $(GED_NATS)

run-agent: $(DIRS)
	go run $(call FIXPATH,./cmd/ged-journal) agent -l debug -f var/agent-log.jsonl --nats $(GED_NATS)


.PHONY: build
build: $(CMDS)

.PHONY: $(CMDS)
$(CMDS): $(DIRS) 
	go build -o $(call FIXPATH,dist/$@$(GOEXE)) $(call FIXPATH,./cmd/$@)

.PHONY: release
release: dist-clean $(DIRS) release_$(GOOS)

release_windows: build
	zip -j $(NAME)_win_$(GOARCH).zip  dist/*

release_linux: build
	cd dist ; tar -czf ../$(NAME)_linux_$(GOARCH).tar.gz *


.PHONY: no-dirty
no-dirty:
	git diff --exit-code

.PHONY: dist-clean
dist-clean:
	$(RM) dist

.PHONY: clean
clean: dist-clean
	$(RM) *.tar.gz
	$(RM) *.zip
