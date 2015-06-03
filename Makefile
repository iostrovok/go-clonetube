GP := $(shell dirname $(realpath $(lastword $(GOPATH))))
ROOT := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
export GOPATH := ${ROOT}/clonetube/:${ROOT}:${GOPATH}


test:
	echo ${ROOT}
	echo ${GOPATH}
	go test ./clonetube/

test-cover:
	echo ${ROOT}
	echo ${GOPATH}
	go test ./clonetube/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out
