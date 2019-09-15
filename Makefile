
dev-dep:
	dep ensure -v


ifndef TARGET
override TARGET = "dev-builder"
endif

ifndef VERSION
override VERSION = "0.0.1-internal"
endif

ifndef
override IMAGE = "kinecosystem/horizon"
endif

DATE := "$(shell date +'%y.%m.%d-%H.%M')"
HOST_MOUNT_POINT := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

build:
	@$(MAKE) tests_teardown
	export HOST_MOUNT_POINT=$(HOST_MOUNT_POINT); \
     	export DATE=$(DATE); \
        export VERSION=$(VERSION);\
        export TARGET=$(TARGET);\
	docker build -f support/images/horizon/Dockerfile \
	--target=$(TARGET) \
        --build-arg DATE=$(DATE) \
	--build-arg VERSION=$(VERSION) \
	--cache-from  kinecosystem/horizon-base:latest \
	-t $(IMAGE)-$(TARGET):latest .

get_horizon:
	export TARGET=$(TARGET);\
	docker rm -f builder || true ; \
	docker run -d --name builder -it kinecosystem/horizon-$(TARGET):latest bash && \
	docker cp builder://go/src/github.com/kinecosystem/go/services/horizon/horizon services/horizon && \
	docker kill builder


#MOUNT_POINT="/jenkins_home/workspace/horizon/go/src/github.com/kinecosystem/go/"
test:
	@$(MAKE) tests_teardown
	docker-compose -f support/images/horizon/docker-compose.yml up -d postgresql mysql redis \
                && export HOST_MOUNT_POINT=$(HOST_MOUNT_POINT); \
                docker-compose -f support/images/horizon/docker-compose.yml run --no-deps  horizon \
                bash -c \
                "dep ensure -v; \
                go get github.com/tebeka/go2xunit; \
                 ./support/scripts/run_tests"
	@$(MAKE) tests_teardown

docker_release: 
	export VERSION=$(VERSION);\
	docker build \
	--target production \
	-f support/images/horizon/Dockerfile \
	--cache-from  kinecosystem/horizon-base:latest \
	 -t $(IMAGE):$(VERSION) \
	 -t $(IMAGE):latest .

docker_push:
	docker push $(IMAGE):$(VERSION)
	#docker push $(IMAGE):latest . || true

tests_teardown:
	export HOST_MOUNT_POINT=$(HOST_MOUNT_POINT);
	docker-compose -f support/images/horizon/docker-compose.yml down -v


jenkins_teardown:
	export HOST_MOUNT_POINT=$(HOST_MOUNT_POINT); \
	docker-compose -f support/images/horizon/docker-compose.yml run --no-deps horizon \
                bash -c \
                "rm -rf Gopkg.lock vendor cover.out test-results.xml services/horizon/horizon"

horizon: 
	@cd services/horizon; ./horizon ${CMD}

lint:
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done


static: vet lint
	go build -i -v -o ${OUT}-v${VERSION} -tags netgo -ldflags="-extldflags \"-static\" -w -s -X main.version=${VERSION}" ${PKG}


dev_clean_go:
	go clean -cache -modcache -i -r

.PHONY: dev_clean_go dev-dep run build test test_setup tests_teardown docker_push horizon jenkins_teardown get_horizon docker_release
