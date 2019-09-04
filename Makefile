
dep:
	dep ensure -v

build:
	docker build --build-arg DATE="1-9-19" --build-arg VERSION="1.2.3" --target=builder -f services/horizon/Dockerfile .

test: 
	@$(MAKE) test_teardown
	docker-compose -f support/images/horizon/docker-compose.yml up -d postgresql mysql redis \
		&& docker-compose -f support/images/horizon/docker-compose.yml run --no-deps horizon ./support/scripts/run_tests
	@$(MAKE) test_teardown

test_teardown:
	docker-compose -f support/images/horizon/docker-compose.yml down -v \
		&& rm -rf support/images/horizon/volumes

test_xunit: test
	go tool cover -func=cover.out

horizon: 
	@cd services/horizon; ./horizon ${CMD}

lint:
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done

static: vet lint
	go build -i -v -o ${OUT}-v${VERSION} -tags netgo -ldflags="-extldflags \"-static\" -w -s -X main.version=${VERSION}" ${PKG}


docker_tag:

docker_push:


clean:
	go clean -cache -modcache -i -r

.PHONY: clean dep run build test test_setup test_teardown test_xunit
