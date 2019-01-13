FROM golang:latest

# Dockerfile for testing only (non Linux OS-es)

# Copy Files to the container


COPY . /go/src/github.com/LinuxArmor/SaFSt-Go

# Install dependencies

RUN go get -d -v github.com/hanwen/go-fuse/fuse

# Add SSH support

EXPOSE 22

# Install FUSE


## Install FUSE
#
#WORKDIR /opt
#
#RUN git clone https://github.com/libfuse/libfuse
#
#WORKDIR /opt/libfuse
#
#RUN mkdir build
#
#WORKDIR /opt/libfuse/build
#
#RUN meson ..
#
#RUN ninja
#
#RUN sudo ninja install

# Add folder to use as an entrypoint

WORKDIR /home
RUN mkdir entrypoint

# Run the project

RUN go run github.com/LinuxArmor/SaFSt-Go/main /home/entrypoint