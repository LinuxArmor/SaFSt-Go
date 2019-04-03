FROM golang:latest

# Dockerfile for testing only (non Linux OS-es)

# Install dependencies

WORKDIR /tmp

RUN apt-get update && \
    apt-get install -y fuse

# Copy Files to the container


COPY . /go/src/github.com/LinuxArmor/SaFSt-Go

# Install dependencies

RUN go get -d -v bazil.org/fuse

# Add SSH support

EXPOSE 22

# Add folder to use as an entrypoint

WORKDIR /home

RUN mkdir entrypoint

# Run the project

RUN go run github.com/LinuxArmor/SaFSt-Go/main /home/entrypoint
