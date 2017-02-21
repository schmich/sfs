FROM golang:latest
MAINTAINER Chris Schmich <schmch@gmail.com>
COPY . /go/src/github.com/schmich/sfs
WORKDIR /go/src/github.com/schmich/sfs
CMD ["/bin/bash", "-c", "/go/src/github.com/schmich/sfs/build-linux.sh"]
