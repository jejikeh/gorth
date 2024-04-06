run: main.go
	@gofmt -w . && go run . run -i=tests/foo.gorth
	
build: main.go
	@gofmt -w . && go run . build -i=tests/foo.gorth -o=out/foo.ssa && qbe -o ./out/foo.s ./out/foo.ssa && cc ./out/foo.s -o ./out/foo && ./out/foo

f:
	@gofmt -w .
	
test:
	go test ./...

qbe:
	@qbe -o ./out/foo.s ./out/foo.ssa && cc ./out/foo.s -o ./out/foo && ./out/foo