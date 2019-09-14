
dev-dep:
	dep ensure -v

# DATE='1-1-1970' VERSION="4.4.4.4"  make build

IMAGE := "kinecosystem/horizon"
build:
	@$(MAKE) test_teardown
	docker build --build-arg DATE=$(DATE) \
	--build-arg VERSION=$(VERSION) \
	-f services/horizon/Dockerfile -t $(IMAGE):$(VERSION) .


#MOUNT_POINT="/jenkins_home/workspace/horizon/go/src/github.com/kinecosystem/go/"
test: 
	@$(MAKE) test_teardown
	docker-compose -f support/images/horizon/docker-compose.yml up -d postgresql mysql redis \
		&& HOST_MOUNT_POINT=$(MOUNT_POINT); \
		docker-compose -f support/images/horizon/docker-compose.yml run --no-deps horizon \
		bash -c \
		"dep ensure -v; \
		go get github.com/tebeka/go2xunit; \
		 ./support/scripts/run_tests"



test_teardown:
	export HOST_MOUNT_POINT=$(MOUNT_POINT); \
	docker-compose -f support/images/horizon/docker-compose.yml run --no-deps horizon \
		bash -c \
		"ls -ltr; \
		rm -rf *"
	docker-compose -f support/images/horizon/docker-compose.yml down -v \
		&& rm -rf support/images/horizon/volumes



horizon: 
	@cd services/horizon; ./horizon ${CMD}

lint:
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done


static: vet lint
	go build -i -v -o ${OUT}-v${VERSION} -tags netgo -ldflags="-extldflags \"-static\" -w -s -X main.version=${VERSION}" ${PKG}


docker_push:
	docker push ${IMAGE} ${IMAGE}:${VERSION}
	#docker push ${IMAGE}:latest


dev_clean_go:
	go clean -cache -modcache -i -r

.PHONY: dev_clean_go dev-dep run build test test_setup test_teardown docker_push horizon
