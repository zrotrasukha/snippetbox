run: build
	@./bin/main 

build: 
	@go build -o ./bin/main ./cmd/web/

test: 
	@go test $(ARGS) ./cmb/web/

