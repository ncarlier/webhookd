# webhookd image.
#
# VERSION 0.0.1
#
# BUILD-USING: docker build --rm -t ncarlier/webhookd .

FROM golang:1.3


# Install ssh-keygen
RUN apt-get update && apt-get install -y ssh sudo

# Install the latest version of the docker CLI
RUN curl -L -o /usr/local/bin/docker https://get.docker.io/builds/Linux/x86_64/docker-latest && \
    chmod +x /usr/local/bin/docker

# Install GO application
WORKDIR /go/src/github.com/ncarlier/webhookd
ADD ./src /go/src/github.com/ncarlier/webhookd
RUN go get github.com/ncarlier/webhookd

# Add scripts
ADD ./scripts /var/opt/webhookd/scripts

# Create work and ssh directories
RUN mkdir /var/opt/webhookd/work

# Generate SSH deploiment key (should be overwrite by a volume)
RUN ssh-keygen -N "" -f /root/.ssh/id_rsa

# Ignor strict host key checking
RUN echo "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config && \
    echo "Host bitbucket.org\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config

# Change workdir
WORKDIR /var/opt/webhookd

# Port
EXPOSE 8080

CMD []
ENTRYPOINT ["/go/bin/webhookd"]
