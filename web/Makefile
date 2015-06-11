.PHONY: build clean run init

BASE=./
PRODUCT := web.goc

# # Prepend our _vendor directory to the system GOPATH
# # so that import path resolution will prioritize
# # our third party snapshots.
GOPATH := ${PWD}/:${GOPATH}

export GOPATH

default: build

build:
	@echo ${GOPATH}
	mkdir ${BASE}/bin
	go build -o ${BASE}/bin/${PRODUCT} ./src/main.go

clean:
	@rm -rf ${BASE}/bin/${PRODUCT}
	@echo "clean ok"

run: clean build
	${BASE}/${PRODUCT} 80

init:
	go get github.com/xuyu/goredis
	mkdir -p /data/logs/cloudbakgo/errorlogs
	mkdir -p /data/logs/cloudbakgo/datalogs
	mkdir -p /data/apps/cloudbakgo
