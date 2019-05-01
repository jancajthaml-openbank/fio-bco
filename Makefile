ifndef GITHUB_RELEASE_TOKEN
$(warning GITHUB_RELEASE_TOKEN is not set)
endif

META := $(shell git rev-parse --abbrev-ref HEAD 2> /dev/null | sed 's:.*/::')
VERSION := $(shell git fetch --tags --force 2> /dev/null; tags=($$(git tag --sort=-v:refname)) && ([ $${\#tags[@]} -eq 0 ] && echo v0.0.0 || echo $${tags[0]}))

.ONESHELL:

.PHONY: all
all: bootstrap sync test package bbtest

.PHONY: package
package:
	@$(MAKE) bundle-binaries
	@$(MAKE) bundle-debian
	@$(MAKE) bundle-docker

.PHONY: bundle-binaries
bundle-binaries:
	@echo "[info] packaging binaries for linux/amd64"
	@docker-compose run --rm package --arch linux/amd64 --pkg fio-bco-import
	@docker-compose run --rm package --arch linux/amd64 --pkg fio-bco-rest

.PHONY: bundle-debian
bundle-debian:
	@echo "[info] packaging for debian"
	@docker-compose run --rm debian -v $(VERSION)+$(META) --arch amd64

.PHONY: bundle-docker
bundle-docker:
	@echo "[info] packaging for docker"
	@docker build -t openbank/fio-bco:$(VERSION)-$(META) .

.PHONY: bootstrap
bootstrap:
	@docker-compose build --force-rm go

.PHONY: lint
lint:
	@docker-compose run --rm lint --pkg fio-bco-import || :
	@docker-compose run --rm lint --pkg fio-bco-rest || :

.PHONY: sec
sec:
	@docker-compose run --rm sec --pkg fio-bco-import || :
	@docker-compose run --rm sec --pkg fio-bco-rest || :

.PHONY: sync
sync:
	@docker-compose run --rm sync --pkg fio-bco-import
	@docker-compose run --rm sync --pkg fio-bco-rest

.PHONY: update
update:
	@docker-compose run --rm update --pkg fio-bco-import
	@docker-compose run --rm update --pkg fio-bco-rest

.PHONY: test
test:
	@docker-compose run --rm test --pkg fio-bco-import
	@docker-compose run --rm test --pkg fio-bco-rest

.PHONY: release
release:
	@docker-compose run --rm release -v $(VERSION)+$(META) -t ${GITHUB_RELEASE_TOKEN}

.PHONY: bbtest
bbtest:
	@docker-compose build bbtest
	@echo "[info] removing older images if present"
	@(docker rm -f $$(docker ps -a --filter="name=fio_bco_bbtest" -q) &> /dev/null || :)
	@echo "[info] running bbtest image"
	@(docker exec -it $$(\
		docker run -d -ti \
			--name=fio_bco_bbtest \
			-e UNIT_VERSION="$(VERSION)-$(META)" \
			-v /sys/fs/cgroup:/sys/fs/cgroup:ro \
			-v /var/run/docker.sock:/var/run/docker.sock \
      -v /var/lib/docker/containers:/var/lib/docker/containers \
			-v $$(pwd)/bbtest:/opt/bbtest \
			-v $$(pwd)/reports:/reports \
			--privileged=true \
			--security-opt seccomp:unconfined \
		openbankdev/fio_bco_bbtest \
	) rspec --require /opt/bbtest/spec.rb \
		--format documentation \
		--format RspecJunitFormatter \
		--out junit.xml \
		--pattern /opt/bbtest/features/*.feature || :)
	@echo "[info] removing bbtest image"
	@(docker rm -f $$(docker ps -a --filter="name=fio_bco_bbtest" -q) &> /dev/null || :)
