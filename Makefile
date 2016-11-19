GOPATH := ${PWD}/../../../../
export GOPATH

BASE="src/go-ops-user-service"
SRC_FOLDERS="service util cache"

#SRCS = $(shell git ls-files '*.go')
SRCS = main.go process.go db/db.go
PKGS = ./. ./service ./moutildels ./cache

glide:
	go get github.com/Masterminds/glide

deps: glide
	glide install

update: glide
	glide update

build: deps
	go build -ldflags "$(LDFLAGS)" -o aws-elk-reports

install: deps
	go install -ldflags "$(LDFLAGS)"

lint:
	@ go get -v github.com/golang/lint/golint
	$(foreach file,$(SRCS),golint $(file) || exit;)

vet:
	#  @-go get -v golang.org/x/tools/cmd/vet
	$(foreach pkg,$(PKGS),go vet $(pkg);)

fmt:
	gofmt -w $(SRCS)

fmtcheck:
	$(foreach file,$(SRCS),gofmt -d $(file);)

docker:
	docker build --no-cache -t aws-elk-reports:latest .
