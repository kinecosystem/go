
# Image used to build kinecosystem/horizon
# kinecosystem/go repo should be mounted into a container,
# then you can execute the build

FROM golang:1.11-stretch as builder

RUN BUILD_DEPS="build-essential git mercurial postgresql-client mysql-client libsodium-dev make"; \
    apt-get -qq update && apt-get -qq install $BUILD_DEPS

RUN  mv /usr/lib/x86_64-linux-gnu/libsodium.so /usr/lib/x86_64-linux-gnu/libsodium.so.18

RUN curl -sS https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN mkdir -p /go/src/github.com/kinecosystem/go

#RUN /bin/bash -c "cd /go/src; go get github.com/kinecosystem/go"
WORKDIR "/go/src/github.com/kinecosystem/go"

ADD  . /go/src/github.com/kinecosystem/go/

RUN /bin/bash -c "/go/bin/dep ensure -v"
RUN /bin/bash -c "cd /services/horizon; go build ."

#VOLUME "/go/src/github.com/kinecosystem/go/services/horizon/"



FROM ubuntu:16.04

COPY  --from=builder "/go/src/github.com/kinecosystem/go/services/horizon/horizon" /usr/local/bin/

RUN chmod 775 /usr/local/bin/horizon
# horizon configuration is done via environment variables
#
# available configuration at:
# https://github.com/kinecosystem/go/blob/horizon-v0.12.3/services/horizon/main.go#L29
ENV DATABASE_URL= \
    HORIZON_DB_MAX_OPEN_CONNECTIONS= \
    STELLAR_CORE_DATABASE_URL= \
    STELLAR_CORE_URL= \
    LOG_LEVEL= \
    INGEST= \
    PER_HOUR_RATE_LIMIT= \
    NETWORK_PASSPHRASE= \
    FRIENDBOT_SECRET= \
    FRIENDBOT_URL= \
    CURSOR_NAME=


EXPOSE 8000

ENTRYPOINT ["/usr/local/bin/horizon"]
CMD ["serve"]


