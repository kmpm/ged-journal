# check OS and set the correct path separator and commands
ifeq ($(OS),Windows_NT)
	FIXPATH = $(subst /,\,$1)
	RM = del /Q
	MKDIR = mkdir
else
	FIXPATH = $1
	RM = rm -f
	MKDIR = mkdir -p
endif

GOEXE:=$(shell go env GOEXE)
CMDNAME=ged-journal

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

.PHONY: test
test:
	go test -v ./...

.PHONY: run
run: $(DIRS)
	go run $(call FIXPATH,./cmd/$(CMDNAME)) run -l debug -f var/ged-journal-log.jsonl


.PHONY: build
build: $(DIRS)
	go build -o $(call FIXPATH,dist/$(CMDNAME)$(GOEXE)) $(call FIXPATH,./cmd/$(CMDNAME))
