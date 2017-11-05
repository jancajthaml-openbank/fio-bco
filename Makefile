NAME = jancajthaml/fio_bco
VERSION = $$(git rev-parse --abbrev-ref HEAD 2> /dev/null | rev | cut -d/ -f1 | rev)
CORES := $$(getconf _NPROCESSORS_ONLN)

.PHONY: all
all: package bundle authors

.PHONY: package
package:
	docker-compose run --rm package

.PHONY: test
test:
	docker-compose run --rm test

.PHONY: bundle
bundle:
	docker-compose build bundle

.PHONY: run
run:
	docker-compose run \
		--rm \
		--no-deps \
		--service-ports \
		-e TENANT_NAME=$(TENANT_NAME) \
		-e ACCOUNT_IBAN=$(ACCOUNT_IBAN) \
		-e FIO_TOKEN=$(FIO_TOKEN) \
		bundle

.PHONY: authors
authors:
	@git log --format='%aN <%aE>' | sort -fu > AUTHORS