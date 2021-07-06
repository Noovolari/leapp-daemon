run:
	go run main.go bootstrap.go
swagger:
	swagger generate spec -o ./swagger.yaml --scan-models

swagger-serve: swagger
	swagger serve ./swagger.yaml
