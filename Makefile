swagger_version="4.14.0"
version=0.1.0
.PHONY: lb-api cloud-provider-manager docker-image setup_init setup_update setup_destroy

all: cloud-provider-manager lb-api

# cloud-provider-manager
cloud-provider-manager: pkg/generate/openapi.gen.go
	go build -v -o cloud-provider-manager cmd/cloud-provider-manager/main.go

docker-image: Dockerfile
	docker build -t shaardie/lb-api-cloud-provider-manager:$(version) .

cloud-provider-manager/generate/openapi.gen.go: openapi.yaml
	oapi-codegen -generate client -package generate openapi.yaml > cloud-provider-manager/generate/openapi.gen.go

# lb-api
lb-api: pkg/generate/openapi.gen.go cmd/lb-api/dist cmd/lb-api/dist/openapi.yaml
	go get -v github.com/shaardie/lb-api/cmd/lb-api
	go build -v -o lb-api cmd/lb-api/main.go

cmd/lb-api/dist:
	curl -L https://github.com/swagger-api/swagger-ui/archive/refs/tags/v$(swagger_version).tar.gz | tar zxv --strip-components=1 swagger-ui-$(swagger_version)/dist
	mv dist cmd/lb-api/
	sed -i 's/https:\/\/petstore.swagger.io\/v2\/swagger.json/openapi.yaml/' cmd/lb-api/dist/swagger-initializer.js

cmd/lb-api/dist/openapi.yaml: openapi.yaml
	cp openapi.yaml cmd/lb-api/dist/

pkg/generate/openapi.gen.go: openapi.yaml oapi.yaml
	oapi-codegen -config oapi.yaml openapi.yaml

setup_init: all
	cd scripts && vagrant up
	cd scripts && vagrant ssh -c "sudo cat /root/.kube/config" > ../kubeconfig

setup_update: all
	cd scripts && vagrant ssh -c "sudo systemctl stop lb-api"
	cd scripts && vagrant upload ../lb-api /src/lb-api
	cd scripts && vagrant ssh -c "sudo systemctl restart lb-api"

	cd scripts && vagrant ssh -c "sudo systemctl stop cloud-provider-manager"
	cd scripts && vagrant upload ../cloud-provider-manager /src/cloud-provider-manager
	cd scripts && vagrant ssh -c "sudo systemctl restart cloud-provider-manager"

setup_destroy:
	cd scripts && vagrant destroy -f
	rm kubeconfig

setup_ssh:
	cd scripts && vagrant ssh
