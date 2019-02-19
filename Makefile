# Parameters
MAIN_FILE=main.go
MAIN_PATH=cmd/streamer_event_server
BIN_PATH=bin
SERVER_MAC=streamer_event_server_mac
SERVER_LINUX=streamer_event_server_linux

build-mac:
	GOARCH=amd64 GOOS=darwin go build -o $(BIN_PATH)/$(SERVER_MAC) $(MAIN_PATH)/$(MAIN_FILE)

build-linux:
	GOARCH=amd64 GOOS=linux go build -o $(BIN_PATH)/$(SERVER_LINUX) $(MAIN_PATH)/$(MAIN_FILE)

run-mac:
	$(BIN_PATH)/$(SERVER_MAC)

run-linux:
	$(BIN_PATH)/$(SERVER_LINUX)