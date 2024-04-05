run: main.go
	@gofmt -w . && go run . run
	
build: main.go
	@gofmt -w . && go run . build && qbe -o file.s file.ssa && cc file.s && ./a.out

f:
	@gofmt -w .
	
test:
	go test ./...