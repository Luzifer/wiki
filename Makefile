default: generate

generate:
	docker run --rm -ti -v $(CURDIR):$(CURDIR) -w $(CURDIR)/src node:alpine \
		sh -exc "npm ci && node ci/build.mjs && chown -R $(shell id -u) ../frontend node_modules"

publish:
	bash ./ci/build.sh
