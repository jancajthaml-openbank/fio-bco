NAME = openbank/fio-bco
VERSION = $$(git rev-parse --abbrev-ref HEAD 2> /dev/null | sed 's:.*/::')
CORES := $$(getconf _NPROCESSORS_ONLN)

.PHONY: all
all: package bundle

.PHONY: package
package:
	docker-compose -f dev/docker-compose.yml \
		run --rm package

.PHONY: test
test:
	docker-compose -f dev/docker-compose.yml \
		run --rm test

.PHONY: bundle
bundle:
	docker-compose -f dev/docker-compose.yml \
		build artefact

.PHONY: run
run:
	docker-compose -f dev/docker-compose.yml \
		run \
		--rm \
		--no-deps \
		--service-ports \
		-e TENANT_NAME=$(TENANT_NAME) \
		-e FIO_TOKEN=$(FIO_TOKEN) \
		artefact
