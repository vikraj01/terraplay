BINARY_NAME=zephyr
BUILD_CMD=go build -o $(BINARY_NAME) ./cmd/zephyr

build:
	$(BUILD_CMD)

clean:
	rm -f $(BINARY_NAME)

build-linux:
	GOOS=linux GOARCH=amd64 $(BUILD_CMD)

build-macos:
	GOOS=darwin GOARCH=amd64 $(BUILD_CMD)

build-windows:
	GOOS=windows GOARCH=amd64 $(BUILD_CMD)
	mv $(BINARY_NAME) $(BINARY_NAME).exe
