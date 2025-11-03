PREFIX ?= $(HOME)
BINARY_DIR = $(PREFIX)/.local/bin
DATA_DIR = $(PREFIX)/.local/share/homework-collection
SYSTEMD_DIR = $(HOME)/.config/systemd/user
BUILD_DIR = build

.PHONY: all install install-binary install-data install-systemd uninstall clean

all: build install

build: clean
	install -d $(BUILD_DIR)
	go build -o $(BUILD_DIR)/homework-collection ./src

install: install-binary install-data install-systemd

install-binary:
	install -d $(BUILD_DIR)
	install -d $(BINARY_DIR)
	install -m755 $(BUILD_DIR)/homework-collection $(BINARY_DIR)/homework-collection

install-data:
	install -d $(DATA_DIR)
	install students.xlsx $(DATA_DIR)/students.xlsx

	install -d $(DATA_DIR)/www/html
	install www/html/index.html $(DATA_DIR)/www/html/index.html
	install www/html/axios.min.js $(DATA_DIR)/www/html/axios.min.js

install-systemd:
	install -d $(SYSTEMD_DIR)
	install ./packaging/homework-collection.service $(SYSTEMD_DIR)/homework-collection.service

uninstall:
	rm $(BINARY_DIR)/homework-collection
	rm -r $(DATA_DIR)
	rm $(SYSTEMD_DIR)/homework-collection.service

clean:
	rm -rf $(BUILD_DIR)
