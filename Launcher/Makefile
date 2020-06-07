NAME := hextranslateInit

GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get

all: build

ifndef VERBOSE
.SILENT:
endif

build:
	echo \\033[36mCompiling hextranslate\\033[0m
	$(GOBUILD) -o $(NAME) -v
	echo \\033[1\;32m      Success!\\033[0m
