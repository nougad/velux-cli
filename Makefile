build:
	go build

genclient:
	../../../../bin/swagger generate client -f ./swagger.yaml

install-deps:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go get -u github.com/go-openapi/errors
	go get -u github.com/go-openapi/runtime
	go get -u github.com/go-openapi/runtime/client
	go get -u github.com/go-openapi/strfmt
