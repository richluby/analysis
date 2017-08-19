# Simplifies the building process

VERSION="0.0.1"
BUILD=`git rev-parse HEAD`

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD}"

default:
	go install ${LDFLAGS}
	${GOPATH}bin/analysis ${ARGS}

clean:
	go clean

