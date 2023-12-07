.PHONY : default
default:
	go build -o bin/main src/cmd/main.go
	./bin/main

.PHONY : clean
clean:
	rm -rf bin/main

.PHONY : test
test:
	cd src/event/ && go test -v
	cd src/socket && go test -v

.PHONY : mock
mock:
	mockgen net Conn > src/mocks/net/conn_mock.go
	mockgen -source=src/socket/server.go Server > src/mocks/socket/server_mock.go


help:
	@echo "  default"
	@echo "  clean"
	@echo "  test"
	@echo "  mock"
