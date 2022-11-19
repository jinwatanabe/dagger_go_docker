FROM docker:dind

RUN apk update 

# Install Golang
RUN apk add --update --no-cache go git make musl-dev curl
RUN export GOPATH=/root/go
RUN export PATH=${GOPATH}/bin:/usr/local/go/bin:$PATH
RUN export GOBIN=$GOROOT/bin
RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin
RUN export GO111MODULE=on
WORKDIR /go/src