run: main.go
	@gofmt -w . && go run . run -i=tests/simple_instructions.gorth
	
build: main.go
	@gofmt -w . && go run . build -i=tests/simple_instructions.gorth && qbe -o ./out/file.s file.ssa && cc ./out/file.s -o ./out/file && ./out/file

f:
	@gofmt -w .
	
test:
	go test ./...