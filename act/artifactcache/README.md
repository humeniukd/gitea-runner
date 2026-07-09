# Swagger-ui

This dist folder use the swagger-ui-dist tools

https://github.com/swagger-api/swagger-ui

How to generate this swagger-ui dist/ files

Installation instruction here https://petstore.swagger.io/v2/swagger.json

Replace https://petstore.swagger.io/v2/swagger.json by `/swagger.yaml`

Command to proceed

```shell
$ rm -r dist/
$ wget https://github.com/swagger-api/swagger-ui/archive/refs/tags/v5.32.6.tar.gz -O swagger-ui.tar.gz
$ tar xvf swagger-ui.tar.gz && rm swagger-ui.tar.gz
$ cp -r swagger-ui-5.32.6/dist/ .
$ rm -r swagger-ui-5.32.6/
$ sed -i 's|https://petstore.swagger.io/v2/swagger.json|/swagger.yaml|g' dist/swagger-initializer.js
```
