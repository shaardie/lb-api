swagger_version="4.14.0"

build: pkg/generate/openapi.gen.go dist dist/openapi.yaml
	go build -v

get:
	go get

dist:
	curl -L https://github.com/swagger-api/swagger-ui/archive/refs/tags/v$(swagger_version).tar.gz | tar zxv --strip-components=1 swagger-ui-$(swagger_version)/dist
	sed -i 's/https:\/\/petstore.swagger.io\/v2\/swagger.json/openapi.yaml/' dist/swagger-initializer.js

dist/openapi.yaml: openapi.yaml
	cp openapi.yaml dist/

pkg/generate/openapi.gen.go: openapi.yaml
	oapi-codegen -package generate openapi.yaml > pkg/generate/openapi.gen.go

clean:
	rm -rf lb-api
