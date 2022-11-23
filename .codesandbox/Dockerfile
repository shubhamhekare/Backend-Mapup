#
# Go Dockerfile
#
# https://github.com/dockerfile/go
#

# Pull base image.
FROM ubuntu:22.04

RUN apt update && apt install -y curl

# Install Go
RUN \
  mkdir -p /goroot && \
  curl https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz | tar xvzf - -C /goroot --strip-components=1

# Set environment variables.
ENV GOROOT /goroot
ENV GOPATH /gopath
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH
