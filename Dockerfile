FROM golang:latest

# Dockerfile for testing only (non Linux OS-es)

# Install dependencies

WORKDIR /tmp

RUN apt-get update && \
    apt-get install -y fuse kmod

# Copy Files to the container


COPY . /go/src/github.com/LinuxArmor/SaFSt-Go

# Install dependencies

RUN go get -d -v ./...

# Add SSH support

EXPOSE 22

# Add folder to use as an entrypoint

WORKDIR /home

RUN mkdir entrypoint

# Run the project

RUN modprobe fuse && \
    chmod 666 /dev/fuse && \
    chown root:$USER /etc/fuse.conf && \
    go run /go/src/github.com/LinuxArmor/SaFSt-Go/main/start.go /home/entrypoint
