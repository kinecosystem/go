#!/usr/bin/env bash
# run all tests using Docker containers for test dependencies:
# Go development environment, PostgreSQL, MySQL and Redis.

set -e
set -x

images_path=${IMAGES_PATH:-support/images/horizon}
dc_path=$images_path/docker-compose.yml

docker-compose -f $dc_path down -v \
    && rm -rf $images_path/volumes \
    && docker-compose -f $dc_path up -d postgresql mysql redis \
    && docker-compose -f $dc_path run --no-deps horizon ./support/scripts/run_tests
