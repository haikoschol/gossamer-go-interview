test:
	go test -cover -v ./...

bench:
	go test -bench=. -run=XYZ -benchmem ./... > bench_result.txt

.PHONY: test bench
