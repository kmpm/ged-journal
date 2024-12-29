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

DIRS = dist

.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v


.PHONY: run
run: 
	go run $(call FIXPATH,./cmd/ged-journal) run


.PHONY: test
test:
	go test -v ./...



$(DIRS):
	$(MKDIR) $@