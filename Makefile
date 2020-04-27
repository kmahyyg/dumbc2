.PHONY: certgen output certgen_debug

clean:
	rm -rf output/*

certgen:
	CGO_ENABLED=0 go build -trimpath -ldflags '-s -w' -o output/certgen cmd/certgen/main.go

certgen_debug:
	CGO_ENABLED=1 go build -o output/certgen cmd/certgen/main.go

