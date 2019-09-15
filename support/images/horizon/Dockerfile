# image used to build kinecosystem/horizon
# kinecosystem/go repo should be mounted into a container,
# then you can execute the build

FROM golang:1.11-stretch as base

RUN BUILD_DEPS="build-essential git mercurial postgresql-client mysql-client"; \
    apt-get -qq update && apt-get -qq install $BUILD_DEPS

RUN curl -sS https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR $GOPATH/src/github.com/kinecosystem/go

COPY Gopkg.toml Gopkg.lock ./

RUN /bin/bash -c "dep ensure --vendor-only -v"

ADD . .


#FROM base as tester
#RUN bash -c "go get github.com/tebeka/go2xunit; \
#                 ./support/scripts/run_tests"

#Build Horizon for production (Linux)
FROM base as builder

ARG VERSION
ARG DATE

RUN /bin/bash -c "cd services/horizon; \
    CGO_ENABLED=0 \
    go build -ldflags=' \
    -X github.com/kinecosystem/go/support/app.buildTime=$DATE \
    -X github.com/kinecosystem/go/support/app.version=$VERSION'"

#Build Horizon for dev (Mac)
FROM base as dev-builder

ARG VERSION
ARG DATE
RUN /bin/bash -c "cd services/horizon; \
    GOOS=darwin \
    go build -ldflags=' \
    -X github.com/kinecosystem/go/support/app.buildTime=$DATE \
    -X github.com/kinecosystem/go/support/app.version=$VERSION'"


#production
FROM ubuntu:16.04 as production

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
