ifndef GITHUB_RELEASE_TOKEN
$(warning GITHUB_RELEASE_TOKEN is not set)
endif

META := $(shell git rev-parse --abbrev-ref HEAD 2> /dev/null | sed 's:.*/::')
VERSION := $(shell git fetch --tags --force 2> /dev/null; tags=($$(git tag --sort=-v:refname)) && ([ $${\#tags[@]} -eq 0 ] && echo v0.0.0 || echo $${tags[0]}))

.ONESHELL:
.PHONY: arm64
.PHONY: amd64
.PHONY: armhf

.PHONY: all
all: bootstrap sync test package bbtest

.PHONY: package
package:
	@$(MAKE) package-amd64
	@$(MAKE) bundle-docker

.PHONY: package-%
package-%: %
	@$(MAKE) bundle-binaries-$^
	@$(MAKE) bundle-debian-$^

.PHONY: bundle-binaries-%
bundle-binaries-%: %
	@docker-compose run --rm package --arch linux/$^ --pkg fio-bco-rest --output /project/packaging/bin
	@docker-compose run --rm package --arch linux/$^ --pkg fio-bco-import --output /project/packaging/bin

.PHONY: bundle-debian-%
bundle-debian-%: %
	@docker-compose run --rm debian --version $(VERSION)+$(META) --arch $^ --source /project/packaging

.PHONY: bundle-docker
bundle-docker:
	@docker build -t openbank/fio-bco:$(VERSION)-$(META) .

.PHONY: bootstrap
bootstrap:
	@docker-compose build --force-rm go

.PHONY: lint
lint:
	@docker-compose run --rm lint --pkg fio-bco-rest || :
	@docker-compose run --rm lint --pkg fio-bco-import || :

.PHONY: sec
sec:
	@docker-compose run --rm sec --pkg fio-bco-rest || :
	@docker-compose run --rm sec --pkg fio-bco-import || :

.PHONY: sync
sync:
	@docker-compose run --rm sync --pkg fio-bco-rest
	@docker-compose run --rm sync --pkg fio-bco-import

.PHONY: test
test:
	@docker-compose run --rm test --pkg fio-bco-rest --output /project/reports
	@docker-compose run --rm test --pkg fio-bco-import --output /project/reports

.PHONY: release
release:
	@docker-compose run --rm release -v $(VERSION)+$(META) -t ${GITHUB_RELEASE_TOKEN}

.PHONY: bbtest
bbtest:
	@(docker rm -f $$(docker ps -a --filter="name=fio_bco_bbtest_amd64" -q) &> /dev/null || :)
	@docker exec -t $$(\
		docker run -d \
			--cpuset-cpus=1 \
			--name=fio_bco_bbtest_amd64 \
			-e IMAGE_VERSION="$(VERSION)-$(META)" \
			-e UNIT_VERSION="$(VERSION)" \
			-e UNIT_ARCH=amd64 \
			-v /var/run/docker.sock:/var/run/docker.sock:rw \
			-v /var/lib/docker/containers:/var/lib/docker/containers:rw \
			-v /sys/fs/cgroup:/sys/fs/cgroup:ro \
			-v $$(pwd)/bbtest:/opt/app \
			-v $$(pwd)/reports:/tmp/reports \
			-w /opt/app \
		jancajthaml/bbtest:amd64 \
	) python3 /opt/app/main.py
	@(docker rm -f $$(docker ps -a --filter="name=fio_bco_bbtest_amd64" -q) &> /dev/null || :)
