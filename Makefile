swagger_version="4.14.0"

.PHONY: lb-api

# lb-api
lb-api: lb-api/generate/openapi.gen.go cmd/lb-api/dist cmd/lb-api/dist/openapi.yaml
	mkdir -p bin
	go build -v -o bin/lb-api cmd/lb-api/main.go

run: build
	./lb-api

cmd/lb-api/dist:
	curl -L https://github.com/swagger-api/swagger-ui/archive/refs/tags/v$(swagger_version).tar.gz | tar zxv --strip-components=1 swagger-ui-$(swagger_version)/dist
	mv dist cmd/lb-api/
	sed -i 's/https:\/\/petstore.swagger.io\/v2\/swagger.json/openapi.yaml/' cmd/lb-api/dist/swagger-initializer.js

cmd/lb-api/dist/openapi.yaml: openapi.yaml
	cp openapi.yaml dist/

lb-api/generate/openapi.gen.go: openapi.yaml
	oapi-codegen -package generate openapi.yaml > lb-api/generate/openapi.gen.go
