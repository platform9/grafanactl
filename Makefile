BIN := grafanactl

.PHONY: $(BIN)

$(BIN):
	go build
