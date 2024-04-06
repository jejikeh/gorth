run: main.go
	@gofmt -w . && go run . run
	
build: main.go
	@gofmt -w . && go run . build && qbe -o ./out/file.s file.ssa && cc ./out/file.s -o ./out/file && ./out/file

f:
	@gofmt -w .
	
test:
	go test ./...