CMD = .
APP = "ignorama_linux_amd64"

default: build

# Default build: fetch dependencies from the internet
build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-X main.Version=$(VERSION)" -o bin/$(APP) $(CMD)

# Build using local vendor directory
build-local: vendor
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -mod=vendor -ldflags "-X main.Version=$(VERSION)" -o bin/$(APP) $(CMD)

# Download all dependencies and create a vendor directory.
vendor:
	go mod tidy
	go mod vendor

clean:
	rm -rf bin/*
	rm -rf vendor
