EXECUTABLE=searchup
BIN_DIR=bin
INSTALL_DIR=~/bin

# None of my targets are based on files, so all are "phony" (https://stackoverflow.com/a/2145605)
.PHONY: bin clean install test

.DEFAULT_GOAL := bin

test:
	go test -cover

bin: test
	go build -o $(BIN_DIR)/$(EXECUTABLE)

install: bin
	cp $(BIN_DIR)/$(EXECUTABLE) $(INSTALL_DIR)/bin

clean:
	rm -rf $(BIN_DIR)