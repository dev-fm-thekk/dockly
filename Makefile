build:
	go build -o dist/dockly ./cmd/dockly

build-all:
	GOOS=linux GOARCH=amd64 go build -o dist/dockly-linux ./cmd/dockly
	GOOS=darwin GOARCH=arm64 go build -o dist/dockly-mac ./cmd/dockly
	GOOS=windows GOARCH=amd64 go build -o dist/dockly.exe ./cmd/dockly