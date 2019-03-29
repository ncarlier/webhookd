#########################################
# Build stage
#########################################
FROM golang:1.11 AS builder

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

# Repository location
ARG REPOSITORY=github.com/ncarlier

# Artifact name
ARG ARTIFACT=webhookd

# Fix lib dep
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# Install deps
RUN apk add --no-cache git openssh-client jq bash

# Install docker-compose
RUN COMPOSE_VERSION="1.23.2" \
&& apk add --no-cache \
  py-pip \
&& pip install --no-cache-dir \
  docker-compose==${COMPOSE_VERSION}

# Create folder structure
RUN mkdir -p /var/opt/webhookd/scripts /var/opt/webhookd/work

# Install binary and default scripts
COPY --from=builder /go/src/$REPOSITORY/$ARTIFACT/release/$ARTIFACT /usr/local/bin/$ARTIFACT
COPY --from=builder /go/src/$REPOSITORY/$ARTIFACT/scripts /var/opt/webhookd/scripts
COPY docker-entrypoint.sh /

# Define entrypoint
ENTRYPOINT ["/docker-entrypoint.sh"]

# Define command
CMD webhookd

