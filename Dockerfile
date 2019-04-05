FROM golang:latest

# Dockerfile for testing only (non Linux OS-es)

# Install dependencies

WORKDIR /tmp

RUN apt-get update && \
    apt-get install -y fuse kmod

# Copy Files to the container


COPY . /go/src/github.com/LinuxArmor/SaFSt-Go

# Install dependencies

RUN go get -u -v bazil.org/fuse github.com/syndtr/goleveldb/leveldb golang.org/x/net/context

# Add SSH support

EXPOSE 22

# Add folder to use as an entrypoint

WORKDIR /home

RUN mkdir entrypoint

# Run the project

RUN modprobe fuse && go run /go/src/github.com/LinuxArmor/SaFSt-Go/main/start.go /home/entrypoint
