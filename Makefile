CORES := $$(getconf _NPROCESSORS_ONLN)

.PHONY: all
all: package test bundle bbtest

.PHONY: package
package:
	docker-compose -f dev/docker-compose.yml \
		run --rm package

.PHONY: test
test:
	docker-compose -f dev/docker-compose.yml \
		run --rm test

.PHONY: bbtest
bbtest:
	docker rm -f $$(docker-compose -f dev/docker-compose-bbtest.yml ps -q) 2> /dev/null || :
	docker-compose -f dev/docker-compose-bbtest.yml run bbtest
	docker rm -f $$(docker-compose -f dev/docker-compose-bbtest.yml ps -q) 2> /dev/null || :

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
