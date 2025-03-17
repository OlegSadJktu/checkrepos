install:
	go build -o $(GOPATH)/bin/checkrepos main.go
	chmod +x $(GOPATH)/bin/checkrepos



install-local:
	if [ -z $(INSTALL_DIR) ]; then \
		@echo "INSTALL_DIR is not set"; \
		@exit 1; \
	fi

	go build -o $(INSTALL_DIR)/checkrepos main.go
	chmod +x $(INSTALL_DIR)/checkrepos

uninstall:
	rm -f $(INSTALL_DIR)/checkrepos
	rm -f $(GOPATH)/bin/checkrepos
