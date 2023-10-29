build:
	go build -o main-$(GOOS)-$(GOARCH)
all:
	$(MAKE) GOOS=linux GOARCH=amd64 build
	$(MAKE) GOOS=linux GOARCH=arm64 build
