#!/usr/bin/env bash
set -e
set -x
docker-compose down -v \
    && rm -rf volumes \
    && docker-compose up -d postgresql mysql redis \
    && docker-compose run --no-deps horizon ./support/scripts/run_tests
