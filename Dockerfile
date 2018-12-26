FROM golang:latest

# Dockerfile for testing only (non Linux OS-es)

# Copy Files to the container

COPY . /go/src/github.com/LinuxArmor/SaFSt-Go

# Install dependencies

RUN go get -d -v github.com/hanwen/go-fuse/fuse

# Add SSH support

EXPOSE 22

# Run the project

RUN go run github.com/LinuxArmor/SaFSt-Go/main