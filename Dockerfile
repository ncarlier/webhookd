#########################################
# Build stage
#########################################
FROM golang:1.8 AS builder
MAINTAINER Nicolas Carlier <n.carlier@nunux.org>

# Repository location
ARG REPOSITORY=github.com/ncarlier

# Artifact name
ARG ARTIFACT=webhookd

# Copy sources into the container
ADD . /go/src/$REPOSITORY/$ARTIFACT

# Set working directory
WORKDIR /go/src/$REPOSITORY/$ARTIFACT

# Build the binary
RUN make

#########################################
# Distribution stage
#########################################
FROM docker:dind
MAINTAINER Nicolas Carlier <n.carlier@nunux.org>

# Repository location
ARG REPOSITORY=github.com/ncarlier

# Artifact name
ARG ARTIFACT=webhookd

# Fix lib dep
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# Install binary
COPY --from=builder /go/src/$REPOSITORY/$ARTIFACT/release/$ARTIFACT-linux-amd64 /usr/local/bin/$ARTIFACT

