run-server:
	cd cmd/server && go run main.go

run-client:
	cd cmd/client && go run main.go

build-server:
	cd cmd/server && go build -o tchat-server main.go && mv tchat-server /usr/local/bin

build-client:
	cd cmd/client && go build -o tchat main.go && mv tchat /usr/local/bin