default: generate

generate:
	docker run --rm -ti -v $(CURDIR):$(CURDIR) -w $(CURDIR)/src node:12-alpine \
		sh -exc "npx npm@lts ci && npx npm@lts run build && chown -R $(shell id -u) ../frontend node_modules"
	go generate

publish:
	curl -sSLo golang.sh https://raw.githubusercontent.com/Luzifer/github-publish/master/golang.sh
	bash golang.sh
