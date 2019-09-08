
dep:
	dep ensure -v

build:
	@$(MAKE) test_teardown
	docker build --build-arg DATE=$(DATE) --build-arg VERSION=$(VERSION) --target=builder -f services/horizon/Dockerfile -t builder .

test: 
	@$(MAKE) test_teardown
	docker-compose -f support/images/horizon/docker-compose.yml up -d postgresql mysql redis \
		&& docker-compose -f support/images/horizon/docker-compose.yml run --no-deps horizon bash -c "go get github.com/tebeka/go2xunit; ./support/scripts/run_tests"


dev-test:
	@$(MAKE) test_teardown
	docker-compose -f support/images/horizon/docker-compose.yml up -d postgresql mysql redis \
		&& docker-compose -f support/images/horizon/docker-compose.yml run --no-deps dev-horizon bash -c "./support/scripts/run_tests"

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

dev_clean_go:
	go clean -cache -modcache -i -r

.PHONY: dev_clean_go dep run build test test_setup test_teardown test_xunit dev-test
