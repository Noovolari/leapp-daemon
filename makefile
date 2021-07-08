run:
	go run main.go bootstrap.go

build_linux:
	echo "> Compiling for Linux"
	GOOS=linux GOARCH=amd64 go build -o bin/leapp-daemon-linux-amd64 main.go bootstrap.go

build_macos:
	echo "> Compiling for MacOS"
	GOOS=darwin GOARCH=amd64 go build -o bin/leapp-daemon-macos-amd64 main.go bootstrap.go

build_windows:
	echo "> Compiling for Windows"
	GOOS=windows GOARCH=amd64 go build -o bin/leapp-daemon-macos-amd64 main.go bootstrap.go

buildall: build_linux build_macos build_windows

swagger:
	swagger generate spec -o ./swagger.yaml --scan-models

swagger-serve: swagger
	swagger serve ./swagger.yaml
